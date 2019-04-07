package epaxos

import (
	"fmt"

	"github.com/bbengfort/epaxos/pb"
)

func (r *Replica) onProposeRequest(e Event) (err error) {
	// Unpack the request from the event
	req := e.Value().(*pb.ProposeRequest)

	// Associate the operation with the source to reply to the client on execute
	r.nops++
	req.Op.Request = r.nops
	source := e.Source().(chan *pb.ProposeReply)
	r.clients[r.nops] = source

	// Create an Instance with the operation
	// QUESTION: What happens to the memory associated with the request? Is it released?
	var inst *pb.Instance
	if inst, err = r.logs.Create(r.PID, []*pb.Operation{req.Op}); err != nil {
		return err
	}

	// Broadcast PreAccept Request for instance
	r.Broadcast(pb.WrapPreacceptRequest(r.Name, &pb.PreacceptRequest{Inst: inst}), false)
	return nil
}

func (r *Replica) onPreacceptRequest(e Event) (err error) {
	// Unpack the request from the event and get the replica log
	req := e.Value().(*pb.PreacceptRequest)
	rlog := r.logs.logs[req.Inst.Replica]

	// Ensure we are appending the next instance (no out of order instances)
	// TODO: move this to the log to return an error if needed.
	if req.Inst.Slot != rlog.nextSlot() {
		return fmt.Errorf("cannot append instance with slot %d into log at slot %d", req.Inst.Slot, rlog.nextSlot())
	}

	// Append the instance into the log for that replica
	rlog.instances = append(rlog.instances, req.Inst)
	changed := r.logs.updateDependencies(req.Inst)
	r.logs.updateConflicts(req.Inst)

	// Prepare the reply
	// TODO: make the channel directional
	source := e.Source().(chan *pb.PeerReply)
	source <- pb.WrapPreacceptReply(r.Name, &pb.PreacceptReply{
		Slot:    req.Inst.Slot,
		Seq:     req.Inst.Seq,
		Deps:    req.Inst.Deps,
		Changed: changed,
	})
	return nil
}

func (r *Replica) onPreacceptReply(e Event) (err error) {
	// Unpack the reply from the event and fetch the instance
	// TODO make the instance fetch from the logs much nicer
	rep := e.Value().(*pb.PreacceptReply)
	inst := r.logs.logs[r.PID].instances[rep.Slot]

	// Ensure the sequence is monotonically increasing
	// TODO: move this helper to the log
	if rep.Seq > r.logs.sequence {
		r.logs.sequence = rep.Seq
	}

	// Only update votes if we're still trying for preaccpt
	if inst.Status == pb.Status_INITIAL {
		inst.Acks++

		// Update dependencies from the remote
		// TODO: Pete does not do this if the quorum == size 3 (optimization?)
		if rep.Changed {
			inst.Changed = true

			// Update the deps and the sequence number
			// TODO: move this functionality to the log
			for pid, dep := range rep.Deps {
				inst.Deps[pid] = dep
			}

			if rep.Seq > inst.Seq {
				inst.Seq = rep.Seq
			}
		}
	}

	// Determine if we've reached a quorum
	// TODO: add the quorum size to the logs
	Q := uint32((len(r.logs.logs) / 2) + 1)

	if inst.Status == pb.Status_INITIAL && inst.Acks >= Q {
		inst.Status = pb.Status_PREACCEPTED
		inst.Acks = 0

		if inst.Changed {
			// Slow Path
			inst.Acks = 1
			r.Broadcast(pb.WrapAcceptRequest(r.Name, &pb.AcceptRequest{Inst: inst}), false)
		} else {
			// Fast Path
			r.Commit(inst)
		}

	}

	return nil
}

func (r *Replica) onBeaconRequest(e Event) (err error) {
	source := e.Source().(chan *pb.PeerReply)
	source <- pb.WrapBeaconReply(r.Name, &pb.BeaconReply{
		QuorumMember: true,
		Replica:      r.PID,
	})

	return nil
}

func (r *Replica) onBeaconReply(e Event) (err error) {
	return nil
}
