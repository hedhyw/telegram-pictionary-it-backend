package game

import (
	"context"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
)

type Controller struct {
	model asyncmodel.RequestEventEmitter
}

func newController(model asyncmodel.RequestEventEmitter) *Controller {
	return &Controller{
		model: model,
	}
}

func (g *Controller) AddPlayer(ctx context.Context, clientID string, meta *telegram.InitDataMeta) error {
	return g.model.EmitRequest(ctx, &RequestEventPlayerJoined{
		ClientID: clientID,
		Meta:     meta,
	})
}

func (g *Controller) RemovePlayer(ctx context.Context, clientID string) error {
	return g.model.EmitRequest(ctx, &RequestEventPlayerRemoved{
		ClientID: clientID,
	})
}

func (g *Controller) Start(ctx context.Context) error {
	return g.model.EmitRequest(ctx, &RequestEventGameStarted{})
}

func (g *Controller) ChangeCanvas(ctx context.Context, clientID string, imageBase64 string) error {
	return g.model.EmitRequest(ctx, &RequestEventCanvasChanged{
		ClientID:    clientID,
		ImageBase64: imageBase64,
	})
}

func (g *Controller) GuessWord(ctx context.Context, clientID string, word string) error {
	return g.model.EmitRequest(ctx, &RequestEventWordGuessAttempted{
		ClientID: clientID,
		Word:     word,
	})
}
