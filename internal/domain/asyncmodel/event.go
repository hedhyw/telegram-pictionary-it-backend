package asyncmodel

import (
	"context"
	"fmt"
)

// Event represent any request/response event.
type Event interface {
	fmt.Stringer
}

// ResponseEvent is an event to a client.
type ResponseEvent interface {
	Event

	// TargetClientIDs returns client ids who should receive this event.
	TargetClientIDs() []string

	// IsResponseEvent is no-op method that indicates response events.
	IsResponseEvent()
}

// RequestEvent is an internal event that should be handled by states.
type RequestEvent interface {
	Event

	// IsRequestEvent is no-op method that indicates request events.
	IsRequestEvent()
}

// RequestEventEmitter creates request events.
type RequestEventEmitter interface {
	// EmitRequest sends the event to the state of the model.
	EmitRequest(ctx context.Context, event RequestEvent) error
}

// ResponseEventConsumer consumes response events.
type ResponseEventConsumer interface {
	// ResponseEvents returns a channel with events, that
	// should be handled by a view.
	ResponseEvents() <-chan ResponseEvent
}
