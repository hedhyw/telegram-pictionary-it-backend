package asyncmodel

import (
	"context"
	"encoding"
	"fmt"
)

// State is a condition of th model in at a specific time.
type State interface {
	HandleRequestEvent(ctx context.Context, event RequestEvent) error

	fmt.Stringer
	encoding.TextMarshaler
}
