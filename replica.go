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

	// Network Definition
	config *Config
	events chan Event
	// remotes map[string]*Remote
	// clients map[uint64]chan *pb.ProposeReply
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
	case ErrorEvent:
		return e.Value().(error)
	default:
		return fmt.Errorf("no handler identified for event %s", e.Type())
	}
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
