package epaxos

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
)

// Names of event types
var eventTypeStrings = [...]string{
	"unknown", "error", "messageReceived", "propose",
	"preacceptRequested", "preacceptReplied", "acceptRequested", "acceptReplied",
	"commitRequested",
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
