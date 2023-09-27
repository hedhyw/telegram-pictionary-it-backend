package asyncmodel

import (
	"context"
	"fmt"
)

const (
	eventsChannelSize = 1024
)

type Model struct {
	currentState State

	handleRequestErrorFunc RequestErrorHandlerFunc

	requestEventsCh  chan RequestEvent
	responseEventsCh chan ResponseEvent
}

func New(
	initialState State,
	handleRequestErrorFunc RequestErrorHandlerFunc,
) *Model {
	model := &Model{
		currentState: initialState,

		handleRequestErrorFunc: handleRequestErrorFunc,

		requestEventsCh:  make(chan RequestEvent, eventsChannelSize),
		responseEventsCh: make(chan ResponseEvent, eventsChannelSize),
	}

	go model.startEventsProcessing(context.TODO())

	return model
}

func (m *Model) SetState(s State) {
	m.currentState = s
}

// nolint: ireturn // State may have different implementations.
func (m Model) State() State {
	return m.currentState
}

func (m Model) ResponseEvents() <-chan ResponseEvent {
	return m.responseEventsCh
}

func (m *Model) startEventsProcessing(ctx context.Context) {
	for event := range m.requestEventsCh {
		if err := m.currentState.HandleRequestEvent(ctx, event); err != nil {
			m.handleRequestErrorFunc(fmt.Errorf("%s: %w", m.currentState, err), event)
		}
	}
}

func (m *Model) EmitResponses(ctx context.Context, events ...ResponseEvent) error {
	for _, e := range events {
		select {
		case m.responseEventsCh <- e:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (m *Model) EmitRequest(ctx context.Context, event RequestEvent) error {
	select {
	case m.requestEventsCh <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
