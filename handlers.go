package epaxos

import "github.com/bbengfort/epaxos/pb"

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
