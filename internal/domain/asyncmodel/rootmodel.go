package asyncmodel

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
)

type Model struct {
	currentState State

	essentials Essentials

	requestEventsCh  chan RequestEvent
	responseEventsCh chan ResponseEvent
}

type Essentials struct {
	InitialState           State
	HandleRequestErrorFunc RequestErrorHandlerFunc
	RequestTimeout         time.Duration
	Logger                 zerolog.Logger
	ChannelSize            int
}

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
