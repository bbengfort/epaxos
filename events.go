package epaxos

import (
	"github.com/bbengfort/epaxos/pb"
)

// Event types represented in ePaxos
const (
	UnknownEvent EventType = iota
	ErrorEvent
	MessageEvent
	ProposeRequestEvent
	PreacceptRequestEvent
	PreacceptReplyEvent
	AcceptRequestEvent
	AcceptReplyEvent
	CommitRequestEvent
	CommitReplyEvent
	BeaconRequestEvent
	BeaconReplyEvent
)

// Names of event types
var eventTypeStrings = [...]string{
	"unknown", "error", "messageReceived", "propose",
	"preacceptRequested", "preacceptReplied", "acceptRequested", "acceptReplied",
	"commitRequested", "commitReplied", "beaconRequested", "beaconReplied",
}

//===========================================================================
// Event Types
//===========================================================================

// EventType is an enumeration of the kind of events that can occur.
type EventType uint16

// String returns the name of event types
func (t EventType) String() string {
	if int(t) < len(eventTypeStrings) {
		return eventTypeStrings[t]
	}
	return eventTypeStrings[0]
}

// Callback is a function that can receive events.
type Callback func(Event) error

//===========================================================================
// Event Definition and Methods
//===========================================================================

// Event represents actions that occur during consensus. Listeners can
// register callbacks with event handlers for specific event types.
type Event interface {
	Type() EventType
	Source() interface{}
	Value() interface{}
}

// event is an internal implementation of the Event interface.
type event struct {
	etype  EventType
	source interface{}
	value  interface{}
}

// Type returns the event type.
func (e *event) Type() EventType {
	return e.etype
}

// Source returns the entity that dispatched the event.
func (e *event) Source() interface{} {
	return e.source
}

// Value returns the current value associated with teh event.
func (e *event) Value() interface{} {
	return e.value
}

//===========================================================================
// Message Event Handlers
//===========================================================================

func requestEvent(req *pb.PeerRequest) *event {
	switch req.Type {
	case pb.Type_PREACCEPT:
		return &event{etype: PreacceptRequestEvent, value: req.GetPreaccept()}
	case pb.Type_ACCEPT:
		return &event{etype: AcceptRequestEvent, value: req.GetAccept()}
	case pb.Type_COMMIT:
		return &event{etype: CommitRequestEvent, value: req.GetCommit()}
	case pb.Type_BEACON:
		return &event{etype: BeaconRequestEvent, value: req.GetBeacon()}
	case pb.Type_UNKNOWN:
		return &event{etype: UnknownEvent, value: req}
	default:
		return &event{etype: UnknownEvent}
	}
}

func replyEvent(rep *pb.PeerReply) *event {
	switch rep.Type {
	case pb.Type_PREACCEPT:
		return &event{etype: PreacceptReplyEvent, value: rep.GetPreaccept()}
	case pb.Type_ACCEPT:
		return &event{etype: AcceptReplyEvent, value: rep.GetAccept()}
	case pb.Type_COMMIT:
		return &event{etype: CommitReplyEvent, value: rep.GetCommit()}
	case pb.Type_BEACON:
		return &event{etype: BeaconReplyEvent, value: rep.GetBeacon()}
	case pb.Type_UNKNOWN:
		return &event{etype: UnknownEvent, value: rep}
	default:
		return &event{etype: UnknownEvent}
	}
}
