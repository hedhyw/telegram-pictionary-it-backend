package asyncmodel

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

// Model is a general model that helps to manage states.
type Model struct {
	currentState State

	essentials Essentials

	requestEventsCh  chan RequestEvent
	responseEventsCh chan ResponseEvent
}

// Essentials contains the required arguments for New.
type Essentials struct {
	InitialState           State
	HandleRequestErrorFunc RequestErrorHandlerFunc
	RequestTimeout         time.Duration
	Logger                 zerolog.Logger
	ChannelSize            int
}

// New creates a new general *Model.
func New(es Essentials) *Model {
	model := &Model{
		currentState: es.InitialState,
		essentials:   es,

		requestEventsCh:  make(chan RequestEvent, es.ChannelSize),
		responseEventsCh: make(chan ResponseEvent, es.ChannelSize),
	}

	go model.startEventsProcessing(context.Background())

	return model
}

// SetState modifies the current state of the model.
func (m *Model) SetState(s State) {
	m.currentState = s
}

// State returns the current state of the event.
//
// nolint: ireturn // State may have different implementations.
func (m *Model) State() State {
	return m.currentState
}

// ResponseEvents returns a channel with events, that
// should be handled by a view.
func (m *Model) ResponseEvents() <-chan ResponseEvent {
	return m.responseEventsCh
}

func (m *Model) startEventsProcessing(ctx context.Context) {
	for event := range m.requestEventsCh {
		m.handleEvent(ctx, event)
	}
}

func (m *Model) handleEvent(ctx context.Context, event RequestEvent) {
	ctx, cancel := context.WithTimeout(ctx, m.essentials.RequestTimeout)
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			m.essentials.Logger.Error().Msgf("panic %v", r)
			debug.PrintStack()
		}
	}()

	if err := m.currentState.HandleRequestEvent(ctx, event); err != nil {
		m.essentials.HandleRequestErrorFunc(fmt.Errorf("%s: %w", m.currentState, err), event)
	}
}

// EmitResponses sends responses to clients.
func (m *Model) EmitResponses(ctx context.Context, events ...ResponseEvent) error {
	ctx, cancel := context.WithTimeout(ctx, m.essentials.RequestTimeout)
	defer cancel()

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

// EmitRequest sends the event to the state of the model.
func (m *Model) EmitRequest(ctx context.Context, event RequestEvent) error {
	ctx, cancel := context.WithTimeout(ctx, m.essentials.RequestTimeout)
	defer cancel()

	select {
	case m.requestEventsCh <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
