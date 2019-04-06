package epaxos

import (
	"context"
	"fmt"
	"io"

	"github.com/bbengfort/epaxos/pb"
)

// Propose is the primary entry point for client requests. This method is the gRPC
// handler that essentially dispatches the propose event to the replica and listens for
// the replica to send back a response so the client can be replied to.
func (r *Replica) Propose(ctx context.Context, in *pb.ProposeRequest) (*pb.ProposeReply, error) {
	// Record the request
	// TODO: add metrics measurements here
	// go func() { r.Metrics.Request(in.Identity) }()

	// Create a channel to wait for the commit handler
	source := make(chan *pb.ProposeReply, 1)

	// Dispatch the event and wait for it to be handled
	event := &event{etype: ProposeRequestEvent, source: source, value: in}
	if err := r.Dispatch(event); err != nil {
		return nil, err
	}

	out := <-source
	return out, nil
}

// Consensus receives PeerRequest messages from remote peers and dispatches them to the
// primary replica process. This method waits for the handler to create a reply before
// receiving the next message.
func (r *Replica) Consensus(stream pb.Epaxos_ConsensusServer) (err error) {
	// currently connected remote peer for logging
	var peer string

	defer func() {
		if peer != "" {
			info("%s disconnected", peer)
		}
	}()

	// Continuously receive messages on the stream
	for {
		var in *pb.PeerRequest
		if in, err = stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if peer == "" {
			peer = in.Sender
			info("%s connected", peer)
		}

		// Unwrap the message and create the specific event type
		e := requestEvent(in)
		if e.Type() == UnknownEvent {
			return fmt.Errorf("received unknown message type from %s", in.Sender)
		}

		// Create the source to wait for the reply
		source := make(chan *pb.PeerReply, 1)
		e.source = source

		// Dispatch the event to the serialized event handler
		// backpressure from this channel will prevent more RECV
		if err = r.Dispatch(e); err != nil {
			return err
		}

		// Wait for the event to be handled before receiving the next
		// message on the stream; this ensures that the order of messages
		// received matches the order of replies sent.
		out := <-source
		if err = stream.Send(out); err != nil {
			return err
		}
	}
}
