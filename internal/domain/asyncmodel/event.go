package asyncmodel

import (
	"context"
	"fmt"
)

type Event interface {
	fmt.Stringer
}

type ResponseEvent interface {
	Event

	TargetClientIDs() []string

	IsResponseEvent()
}

type RequestEvent interface {
	Event

	IsRequestEvent()
}

type RequestEventEmitter interface {
	EmitRequest(ctx context.Context, event RequestEvent) error
}

type ResponseEventConsumer interface {
	ResponseEvents() <-chan ResponseEvent
}
