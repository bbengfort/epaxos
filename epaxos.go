/*
Package epaxos implements the ePaxos consensus algorithm.
*/
package epaxos

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bbengfort/epaxos/pb"
	"github.com/bbengfort/x/noplog"
	"github.com/bbengfort/x/peers"
	"google.golang.org/grpc/grpclog"
)

//===========================================================================
// Package Initialization
//===========================================================================

// PackageVersion of the current ePaxos implementation
const PackageVersion = "0.1"

// Initialize the package and random numbers, etc.
func init() {
	// Set the random seed to something different each time.
	rand.Seed(time.Now().UnixNano())

	// Initialize our debug logging with our prefix
	SetLogger(log.New(os.Stdout, "[epaxos] ", log.Lmicroseconds))
	cautionCounter = new(counter)
	cautionCounter.init()

	// Stop the grpc verbose logging
	grpclog.SetLogger(noplog.New())
}

//===========================================================================
// New ePaxos Instance
//===========================================================================

// New ePaxos replica with the specified config.
func New(options *Config) (replica *Replica, err error) {
	// Create a new configuration from defaults, configuration file, and
	// the environment; then verify it, returning any errors.
	config := new(Config)
	if err = config.Load(); err != nil {
		return nil, err
	}

	// Update the configuration with the passed in configuration.
	if err = config.Update(options); err != nil {
		return nil, err
	}

	// Set the logging level and the random seed
	SetLogLevel(uint8(config.LogLevel))
	if config.Seed != 0 {
		debug("setting random seed to %d", config.Seed)
		rand.Seed(config.Seed)
	}

	// Create and initialize the replica
	replica = new(Replica)
	replica.config = config
	replica.thrifty = config.GetThrifty()
	replica.clients = make(map[uint64]chan *pb.ProposeReply)
	// replica.log = NewLog(replica)
	// replica.Metrics = NewMetrics()

	// Create the local replica definition
	replica.Peer, err = config.GetPeer()
	if err != nil {
		return nil, err
	}

	// Fetch all remote peers (e.g. all peers but self)
	// NOTE: we expect peers to be sorted by PID
	var peers []peers.Peer
	if peers, err = config.GetRemotes(); err != nil {
		return nil, err
	}

	// Make the remotes to send message to remote peers
	replica.remotes = make(Remotes)
	for _, peer := range peers {
		if replica.remotes[peer.PID], err = NewRemote(peer, replica); err != nil {
			return nil, err
		}
	}

	// Set state to initialized
	info("epaxos replica with %d remote peers created", len(replica.remotes))
	if replica.thrifty != nil {
		pids := make([]string, 0, len(replica.thrifty))
		for _, p := range replica.thrifty {
			pids = append(pids, fmt.Sprintf("%d", p))
		}

		info("thrifty communications to replica PIDs %s", strings.Join(pids, ", "))
	}
	return replica, nil
}
