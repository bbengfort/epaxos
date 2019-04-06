package epaxos

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bbengfort/epaxos/pb"
	"github.com/bbengfort/x/peers"
	"google.golang.org/grpc"
)

// DefaultRetries specifies the number of times to attempt a commit.
const DefaultRetries = 3

// NewClient creates a new ePaxos client to connect to a quorum.
func NewClient(remote string, options *Config) (client *Client, err error) {
	// Create a new configuration from defaults, configuration file, and the
	// environment; verify it and return any improperly configured errors.
	config := new(Config)
	if err = config.Load(); err != nil {
		return nil, err
	}

	// Update the configuration with the specified options
	if err = config.Update(options); err != nil {
		return nil, err
	}

	// Create the client
	client = &Client{config: config}

	// Compute the identity
	hostname, _ := config.GetName()
	if hostname != "" {
		client.identity = fmt.Sprintf("%s-%04X", hostname, rand.Intn(0x10000))
	} else {
		client.identity = fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
	}

	// Connect when client is created to capture any errors as early as possible.
	// NOTE: connection errors still return the client for retry
	if err = client.connect(remote); err != nil {
		return client, err
	}

	return client, nil
}

// Client maintains network information embedded in the configuration to
// connect to an ePaxos consensus quorum and make propose requests.
type Client struct {
	sync.RWMutex
	config   *Config          // network details for connection
	conn     *grpc.ClientConn // grpc connection to dial an ePaxos server
	client   pb.EpaxosClient  // grpc RPC interface
	identity string           // a unique identity for all clients
}

//===========================================================================
// Request API
//===========================================================================

// Get a value for a key (execute a read operation)
func (c *Client) Get(key string) ([]byte, error) {
	rep, err := c.Propose(pb.AccessType_READ, key, nil)
	if err != nil {
		return nil, err
	}
	return rep.Value, nil
}

// Put a value for a key (execute a write operation)
func (c *Client) Put(key string, value []byte, execute bool) error {
	var access pb.AccessType
	if execute {
		access = pb.AccessType_WRITEREAD
	} else {
		access = pb.AccessType_WRITE
	}

	_, err := c.Propose(access, key, value)
	return err
}

// Del a value for a key (execute a delete operation)
func (c *Client) Del(key string) error {
	_, err := c.Propose(pb.AccessType_DELETE, key, nil)
	return err
}

// Propose an operation to be applied to the state store.
func (c *Client) Propose(access pb.AccessType, key string, value []byte) (rep *pb.ProposeReply, err error) {
	// Create the request for the operation
	req := &pb.ProposeRequest{
		Identity: c.identity,
		Op: &pb.Operation{
			Type: access, Key: key, Value: value,
		},
	}

	// Send the request
	if rep, err = c.send(req, DefaultRetries); err != nil {
		return nil, err
	}
	return rep, nil
}

// Send the propose request, handling retries.
func (c *Client) send(req *pb.ProposeRequest, retries int) (*pb.ProposeReply, error) {
	// Don't attempt if there are no more retries.
	if retries <= 0 {
		return nil, ErrRetries
	}

	// Connect if not connected
	if !c.isConnected() {
		if err := c.connect(""); err != nil {
			return nil, err
		}
	}

	// Create the context
	timeout, err := c.config.GetTimeout()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rep, err := c.client.Propose(ctx, req)
	if err != nil {
		if retries > 1 {
			// If there is an error connecting to the current host, try another
			if err = c.connect(""); err != nil {
				return nil, err
			}
			return c.send(req, retries-1)
		}
		return nil, err
	}

	if !rep.Success {
		// If there was an error, return the error, otherwise retry.
		if rep.Error != "" {
			return nil, errors.New(rep.Error)
		}
		return c.send(req, retries-1)
	}

	return rep, nil
}

//===========================================================================
// Connection Handlers
//===========================================================================

// Connect to the remote client using the specified timeout. If a remote is not
// specified (e.g. empty string) then a random replica is selected from the
// configuration to connect to, prioritizing any replica on the same host as the
// client.
func (c *Client) connect(remote string) (err error) {
	// Close the connection if one is already open.
	c.close()

	// Get the peer by name or select peer from configuration.
	var host *peers.Peer
	if host, err = c.selectRemote(remote); err != nil {
		return err
	}

	// Parse timeout from configuration
	var timeout time.Duration
	if timeout, err = c.config.GetTimeout(); err != nil {
		return err
	}

	// Connect to the remote's address
	addr := host.Endpoint(false)
	if c.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(timeout)); err != nil {
		return fmt.Errorf("could not connect to '%s': %s", addr, err)
	}

	// Create gRPC client and return
	c.client = pb.NewEpaxosClient(c.conn)
	return nil
}

// Close the connection to the remote host and clean up.
func (c *Client) close() (err error) {
	defer func() {
		c.conn = nil
		c.client = nil
	}()

	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
}

// Find the remote by name, or select the remote on the same host as the client,
// otherwise return a random remote. Future versions should select a remote in the
// same region as the client.
func (c *Client) selectRemote(remote string) (*peers.Peer, error) {
	if remote == "" {
		if len(c.config.Peers) == 0 {
			return nil, ErrNoNetwork
		}

		// Search for host by client name/local hostname
		hostname, _ := c.config.GetName()
		if hostname != "" {
			for _, peer := range c.config.Peers {
				if peer.Name == hostname || peer.Hostname == hostname {
					return &peer, nil
				}
			}
		}

		// Return a random peer
		idx := rand.Intn(len(c.config.Peers))
		return &c.config.Peers[idx], nil
	}

	for _, peer := range c.config.Peers {
		if peer.Name == remote {
			return &peer, nil
		}
	}

	return nil, fmt.Errorf("could not find remote '%s' in configuration", remote)
}

// Ensures a client and connection exist
func (c *Client) isConnected() bool {
	return c.client != nil && c.conn != nil
}
