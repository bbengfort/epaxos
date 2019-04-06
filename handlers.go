package epaxos

import "github.com/bbengfort/epaxos/pb"

func (r *Replica) onProposeRequest(e Event) (err error) {
	// TODO: actually implement propose
	source := e.Source().(chan *pb.ProposeReply)
	source <- &pb.ProposeReply{
		Success: false,
		Error:   "not implemented by this server yet",
	}

	return nil
}

func (r *Replica) onBeaconRequest(e Event) (err error) {
	source := e.Source().(chan *pb.PeerReply)
	source <- pb.WrapBeaconReply(r.Name, &pb.BeaconReply{
		QuorumMember: true,
		Replica:      int64(r.PID),
	})

	return nil
}

func (r *Replica) onBeaconReply(e Event) (err error) {
	return nil
}
