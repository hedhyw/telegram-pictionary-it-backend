package asyncmodel

import (
	"context"
	"fmt"
	"time"
)

const (
	eventsChannelSize = 1024
)

type Model struct {
	currentState   State
	requestTimeout time.Duration

	handleRequestErrorFunc RequestErrorHandlerFunc

	requestEventsCh  chan RequestEvent
	responseEventsCh chan ResponseEvent
}

func New(
	initialState State,
	handleRequestErrorFunc RequestErrorHandlerFunc,
	requestTimeout time.Duration,
) *Model {
	model := &Model{
		currentState:   initialState,
		requestTimeout: requestTimeout,

		handleRequestErrorFunc: handleRequestErrorFunc,

		requestEventsCh:  make(chan RequestEvent, eventsChannelSize),
		responseEventsCh: make(chan ResponseEvent, eventsChannelSize),
	}

	go model.startEventsProcessing(context.Background())

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
		m.handleEvent(ctx, event)
	}
}

func (m *Model) handleEvent(ctx context.Context, event RequestEvent) {
	ctx, cancel := context.WithTimeout(ctx, m.requestTimeout)
	defer cancel()

	if err := m.currentState.HandleRequestEvent(ctx, event); err != nil {
		m.handleRequestErrorFunc(fmt.Errorf("%s: %w", m.currentState, err), event)
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
