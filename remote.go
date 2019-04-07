package epaxos

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/bbengfort/epaxos/pb"
	"github.com/bbengfort/x/peers"
	"google.golang.org/grpc"
)

// MessageBufferSize represents the number of messages that can be queued to send to the
// remote server before the process starts to block to prevent back pressure.
const MessageBufferSize = 1024

// Remote maintains a connection to a peer on the network.
type Remote struct {
	sync.Mutex
	peers.Peer

	sender   string                    // the name of the sender to attach to all messages
	actor    Actor                     // the listener to dispatch events to
	timeout  time.Duration             // timeout before dropping message
	conn     *grpc.ClientConn          // grpc dial connection to the remote
	client   pb.EpaxosClient           // rpc client specified by protobuf
	stream   pb.Epaxos_ConsensusClient // consensus messages stream
	online   bool                      // if the client is connected or not
	messages chan *pb.PeerRequest      // internal channel to schedule messages to be sent on
	done     chan error                // used to wait until the internal go routine is done
}

// Remotes is a collection of remote peers that must be ordered by PID
type Remotes map[uint32]*Remote

//===========================================================================
// External Interface
//===========================================================================

// NewRemote creates a new remote associated with the replica.
func NewRemote(p peers.Peer, r *Replica) (*Remote, error) {
	timeout, err := r.config.GetTimeout()
	if err != nil {
		return nil, err
	}

	remote := &Remote{Peer: p, actor: r, sender: r.Name, timeout: timeout}
	return remote, nil
}

// Connect to the remote peer and establish birdirectional streams. This is the external
// version of connect that is used to establish the go routine that will send and
// receive messages from the remote host and dispatch them to the replic'as event chan.
// It is the messenger routine's responsibility to keep the channel open.
func (c *Remote) Connect() (err error) {
	c.Lock()
	defer c.Unlock()

	if c.messages != nil {
		return fmt.Errorf("already connected to %s (%s)", c.Name, c.Endpoint(true))
	}

	// Establish the messenger channels
	c.messages = make(chan *pb.PeerRequest, MessageBufferSize)
	c.done = make(chan error)

	// Run the messenger go routine.
	go c.messenger()
	return nil
}

// Close the messenger process gracefully by closing the messages channel, wait for the
// messenger to finish sending the last messages, then clean up the connection to the
// remote peer. Note that any sends after Close() will cause a panic.
func (c *Remote) Close() (err error) {
	c.Lock()
	if c.messages == nil {
		c.Unlock()
		return fmt.Errorf("remote messenger to %s (%s) is not running", c.Name, c.Endpoint(true))
	}

	close(c.messages)
	c.Unlock()
	return <-c.done
}

// Send a message to the remote. This places the message on a buffered channel, which
// will be sent in the order they are received. The response is dispatched to the actor
// event listener to be handled in the order responses are received.
func (c *Remote) Send(req *pb.PeerRequest) {
	c.messages <- req
}

//===========================================================================
// RPC Wrappers
//===========================================================================

// SendBeacon sends a beacon message that establishes a conenction to the remote.
func (c *Remote) SendBeacon() {
	beacon := &pb.BeaconRequest{}
	c.Send(pb.WrapBeaconRequest(c.sender, beacon))
}

//===========================================================================
// Messenger Routine
//===========================================================================

// Messenger should run in its own go routine; it sends messages from external threads
// to the remote peer, waiting for a response, which dispatches the reply to the actor
// as an event so that replies are handled in order. There should only be one messenger
// thread running at a time.
func (c *Remote) messenger() {
	// Attempt to establish a connection to the remote peer
	c.connect()

	// Send messages and receive responses
	for msg := range c.messages {
		// If we're not online try to re-establish the connection
		if !c.online {
			if err := c.connect(); err != nil {
				caution("dropped %s message to %s (%s): could not connect", msg.Type, c.Name, c.Endpoint(true))
				c.close()
				continue
			}
		}

		// Send the peer request message
		if err := c.stream.Send(msg); err != nil {
			// go offline if there was an error sending the message
			caution("dropped %s message to %s (%s): could not send", msg.Type, c.Name, c.Endpoint(true))
			c.close()
			continue
		}

		// Wait for the peer reply message
		rep, err := c.stream.Recv()
		if err != nil {
			if err != io.EOF {
				caution(err.Error())
			} else {
				caution("stream to %s (%s) closed by remote", c.Name, c.Endpoint(true))
			}
			c.stream = nil
			c.close()
			continue
		}

		// Dispatch the event to the replica
		if err := c.actor.Dispatch(replyEvent(rep)); err != nil {
			caution("could not dispatch message from %s (%s): %s", c.Name, c.Endpoint(true), err)
		}

	}
}

// Connect to the remote using the specified timeout. Connect is usually not explicitly
// called, but is connected when a message is sent.
func (c *Remote) connect() (err error) {
	// If we're already online do not connect
	if c.online {
		return nil
	}

	addr := c.Endpoint(true)

	if c.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(c.timeout)); err != nil {
		return fmt.Errorf("could not connect to '%s': %s", addr, err)
	}

	// NOTE: do not set online to true until after a response from remote.
	c.client = pb.NewEpaxosClient(c.conn)

	// Create the messages stream
	if c.stream, err = c.client.Consensus(context.Background()); err != nil {
		return fmt.Errorf("could not create peer to peer stream to '%s': %s", addr, err)
	}

	// Always send beacon when connected to establish link.
	c.SendBeacon()

	// mark connection as online
	c.online = true
	return nil
}

// Close the connection to the remote and clean up the connection objects.
func (c *Remote) close() (err error) {
	// Ensure a valid state after close
	defer func() {
		c.conn = nil
		c.client = nil
		c.stream = nil
		c.online = false
	}()

	if c.stream != nil {
		if err = c.stream.CloseSend(); err != nil {
			return fmt.Errorf("could not close stream to %s: %s", c.Endpoint(true), err)
		}
	}

	// Don't cause any panics if already closed
	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			return fmt.Errorf("could not close connection to %s: %s", c.Endpoint(true), err)
		}
	}

	return nil
}

// Reset the connection to the remote in order to try another message
func (c *Remote) reset() (err error) {
	if err = c.close(); err != nil {
		return err
	}
	return c.connect()
}
