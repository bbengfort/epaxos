package epaxos

import (
	"errors"
	"fmt"
	"net"

	"github.com/bbengfort/epaxos/pb"
	"github.com/bbengfort/x/peers"
	"google.golang.org/grpc"
)

// Replica represents the local consensus replica and is the primary object implemented
// in a running system. There should only be one replica per process.
type Replica struct {
	peers.Peer

	config  *Config                          // the static configuration of the replica
	events  chan Event                       // serialize events in the system in the order they're received
	remotes Remotes                          // connections to remote peers to send messages to
	thrifty []uint16                         // the peers to send broadcast messages to
	clients map[uint64]chan *pb.ProposeReply // connected clients awaiting a reply
}

// Listen for messages from peers and clients and run the event loop.
func (r *Replica) Listen() error {
	// Open TCP socket to listen for messages
	addr := fmt.Sprintf(":%d", r.Port)
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s", addr)
	}
	defer sock.Close()
	info("listening for requests on %s", addr)

	// Create the events channel
	r.events = make(chan Event, actorEventBufferSize)

	// Initialize and run the gRPC server in its own thread
	srv := grpc.NewServer()
	pb.RegisterEpaxosServer(srv, r)
	go srv.Serve(sock)

	// Open up connections to remote peers
	if err := r.Connect(); err != nil {
		return err
	}

	// Run the event handling loop
	if r.config.Aggregate {
		if err := r.runAggregatingEventLoop(); err != nil {
			return err
		}
	} else {
		if err := r.runEventLoop(); err != nil {
			return err
		}
	}

	return nil
}

// Close the event handler and stop listening for events.
// TODO: gracefully shutdown the grpc server as well.
func (r *Replica) Close() error {
	if r.events == nil {
		return ErrNotListening
	}
	close(r.events)
	return nil
}

// Dispatch events by clients to the replica.
func (r *Replica) Dispatch(e Event) error {
	if r.events == nil {
		return ErrNotListening
	}

	r.events <- e
	return nil
}

// Handle the events in serial order.
func (r *Replica) Handle(e Event) error {
	trace("%s event received: %v", e.Type(), e.Value())

	switch e.Type() {
	case ProposeRequestEvent:
		return r.onProposeRequest(e)
	case BeaconRequestEvent:
		return r.onBeaconRequest(e)
	case BeaconReplyEvent:
		return r.onBeaconReply(e)
	case ErrorEvent:
		return e.Value().(error)
	default:
		return fmt.Errorf("no handler identified for event %s", e.Type())
	}
}

// Connect the replica to its remote peers. If in thrifty mode, only connects to its
// thrifty neighbors rather than establishing connections to all peers.
func (r *Replica) Connect() error {
	if r.thrifty != nil {
		// Only open up connections to thrifty peers
		for _, pid := range r.thrifty {
			if err := r.remotes[pid].Connect(); err != nil {
				return err
			}
		}
	} else {
		// Open up connections to all peers
		for _, remote := range r.remotes {
			if err := remote.Connect(); err != nil {
				return err
			}
		}
	}

	return nil
}

//===========================================================================
// Event Loops
//===========================================================================

// Runs a normal event loop, handling one event at a time.
func (r *Replica) runEventLoop() error {
	defer func() {
		// nilify the events channel when we stop running it
		r.events = nil
	}()

	for e := range r.events {
		if err := r.Handle(e); err != nil {
			return err
		}
	}

	return nil
}

// Runs an event loop that aggregates multiple propose requests into a single
// instance that is sent to all peers at once.
func (r *Replica) runAggregatingEventLoop() error {
	return errors.New("aggregating operations not implemented yet")
}
