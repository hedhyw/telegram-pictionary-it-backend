package asyncmodel

import (
	"context"
	"encoding"
	"fmt"
)

type State interface {
	HandleRequestEvent(ctx context.Context, event RequestEvent) error

	fmt.Stringer
	encoding.TextMarshaler
}

type StateSetter interface{}
