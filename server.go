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
	// Continuously receive messages on the stream
	for {
		var in *pb.PeerRequest
		if in, err = stream.Recv(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		// Unwrap the message and create the specific event type
		var e *event
		switch in.Type {
		case pb.Type_PREACCEPT:
			e = &event{etype: PreacceptRequestEvent, value: in.GetPreaccept()}
		case pb.Type_ACCEPT:
			e = &event{etype: AcceptRequestEvent, value: in.GetAccept()}
		case pb.Type_COMMIT:
			e = &event{etype: CommitRequestEvent, value: in.GetCommit()}
		default:
			return fmt.Errorf("could not handle message of type %s", in.Type)
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
