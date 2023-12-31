package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/entities"
)

type handlerFunc func(ctx context.Context, clientID string, request entities.SocketRequest) error

// Controller manipulates with the model.
type Controller struct {
	model asyncmodel.RequestEventEmitter

	handlers map[string]handlerFunc
}

func newController(model asyncmodel.RequestEventEmitter) *Controller {
	controller := &Controller{
		model:    model,
		handlers: map[string]handlerFunc{},
	}

	controller.handlers[RequestEventGameStarted{}.String()] = controller.StartGame
	controller.handlers[RequestEventCanvasChanged{}.String()] = controller.ChangeCanvas
	controller.handlers[RequestEventWordGuessAttempted{}.String()] = controller.GuessWord
	controller.handlers[RequestEventPlayerInitialized{}.String()] = controller.InitializePlayer

	return controller
}

// RegisterClient adds a new client to the hub.
func (c *Controller) RegisterClient(
	ctx context.Context,
	clientID string,
	eventsCh chan<- json.RawMessage,
) error {
	return c.model.EmitRequest(ctx, &RequestEventClientConnected{
		ClientID: clientID,
		EventsCh: eventsCh,
	})
}

// UnregisterClient removes the client from the hub.
func (c *Controller) UnregisterClient(
	ctx context.Context,
	clientID string,
) error {
	return c.model.EmitRequest(ctx, &RequestEventClientDisconnnected{
		ClientID: clientID,
	})
}

// EmitClientEvent processes general client request events.
func (c *Controller) EmitClientEvent(
	ctx context.Context,
	clientID string,
	event json.RawMessage,
) error {
	var payload entities.SocketRequest

	err := json.Unmarshal(event, &payload)
	if err != nil {
		return fmt.Errorf("unmarshaling payload: %w", err)
	}

	handler, ok := c.handlers[payload.Name]
	if !ok {
		return semerr.NewNotFoundError(semerr.Error("handler is not found: " + payload.Name))
	}

	return handler(ctx, clientID, payload)
}

// StartGame handles starting of a new game by the client.
func (c *Controller) StartGame(
	ctx context.Context,
	clientID string,
	_ entities.SocketRequest,
) (err error) {
	return c.model.EmitRequest(ctx, &RequestEventGameStarted{
		ClientID: clientID,
	})
}

// ChangeCanvas handles a drawing from the leader.
func (c *Controller) ChangeCanvas(
	ctx context.Context,
	clientID string,
	request entities.SocketRequest,
) (err error) {
	var payload RequestEventCanvasChanged

	err = json.Unmarshal(request.Payload, &payload)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("decoding payload: %w", err))
	}

	payload.ClientID = clientID

	return c.model.EmitRequest(ctx, &payload)
}

// GuessWord handles an attempt to guess the word.
func (c *Controller) GuessWord(
	ctx context.Context,
	clientID string,
	request entities.SocketRequest,
) (err error) {
	var payload RequestEventWordGuessAttempted

	err = json.Unmarshal(request.Payload, &payload)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("decoding payload: %w", err))
	}

	payload.ClientID = clientID

	return c.model.EmitRequest(ctx, &payload)
}

// InitializePlayer handles InitData from Telegram.
func (c *Controller) InitializePlayer(
	ctx context.Context,
	clientID string,
	request entities.SocketRequest,
) (err error) {
	var payload RequestEventPlayerInitialized

	err = json.Unmarshal(request.Payload, &payload)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("decoding payload: %w", err))
	}

	payload.ClientID = clientID

	return c.model.EmitRequest(ctx, &payload)
}
