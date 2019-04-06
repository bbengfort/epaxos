package pb

//===========================================================================
// Message Wrappers
//===========================================================================

// WrapPreacceptRequest in a PeerRequest for transmission on a single streaming channel.
func WrapPreacceptRequest(sender string, msg *PreacceptRequest) *PeerRequest {
	return &PeerRequest{
		Type:   Type_PREACCEPT,
		Sender: sender,
		Message: &PeerRequest_Preaccept{
			Preaccept: msg,
		},
	}
}

// WrapPreacceptReply in a PeerReply for transmission on a single streaming channel.
func WrapPreacceptReply(sender string, msg *PreacceptReply) *PeerReply {
	return &PeerReply{
		Type:    Type_PREACCEPT,
		Sender:  sender,
		Success: true,
		Message: &PeerReply_Preaccept{
			Preaccept: msg,
		},
	}
}

// WrapAcceptRequest in a PeerRequest for transmission on a single streaming channel.
func WrapAcceptRequest(sender string, msg *AcceptRequest) *PeerRequest {
	return &PeerRequest{
		Type:   Type_ACCEPT,
		Sender: sender,
		Message: &PeerRequest_Accept{
			Accept: msg,
		},
	}
}

// WrapAcceptReply in a PeerReply for transmission on a single streaming channel.
func WrapAcceptReply(sender string, msg *AcceptReply) *PeerReply {
	return &PeerReply{
		Type:    Type_ACCEPT,
		Sender:  sender,
		Success: true,
		Message: &PeerReply_Accept{
			Accept: msg,
		},
	}
}

// WrapCommitRequest in a PeerRequest for transmission on a single streaming channel.
func WrapCommitRequest(sender string, msg *CommitRequest) *PeerRequest {
	return &PeerRequest{
		Type:   Type_COMMIT,
		Sender: sender,
		Message: &PeerRequest_Commit{
			Commit: msg,
		},
	}
}

// WrapCommitReply in a PeerReply for transmission on a single streaming channel.
func WrapCommitReply(sender string, msg *CommitReply) *PeerReply {
	return &PeerReply{
		Type:    Type_COMMIT,
		Sender:  sender,
		Success: true,
		Message: &PeerReply_Commit{
			Commit: msg,
		},
	}
}

// WrapBeaconRequest in a PeerRequest for transmission on a single streaming channel.
func WrapBeaconRequest(sender string, msg *BeaconRequest) *PeerRequest {
	return &PeerRequest{
		Type:   Type_BEACON,
		Sender: sender,
		Message: &PeerRequest_Beacon{
			Beacon: msg,
		},
	}
}

// WrapBeaconReply in a PeerReply for transmission on a single streaming channel.
func WrapBeaconReply(sender string, msg *BeaconReply) *PeerReply {
	return &PeerReply{
		Type:    Type_BEACON,
		Sender:  sender,
		Success: true,
		Message: &PeerReply_Beacon{
			Beacon: msg,
		},
	}
}
