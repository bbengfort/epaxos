package epaxos

import "errors"

// Standard errors for primary operations.
var (
	ErrEventTypeError   = errors.New("captured event with wrong value type")
	ErrEventSourceError = errors.New("captured event with wrong source type")
	ErrUnknownState     = errors.New("epaxos in an unknown state")
	ErrNotListening     = errors.New("replica is not listening for events")
	ErrRetries          = errors.New("could not connect after several attempts")
	ErrNoNetwork        = errors.New("no network specified in the configuration")
	ErrBenchmarkMode    = errors.New("specify either fixed duration or maximum operations benchmark mode")
	ErrBenchmarkRun     = errors.New("benchmark has already been run")
)
