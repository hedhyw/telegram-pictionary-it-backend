package game

import (
	"context"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/entities"
)

// TODO: async view.
// TODO: exported function comments.
// TODO: tests.
type view struct {
	essentials Essentials

	model asyncmodel.ResponseEventConsumer
}

func newView(es Essentials, model asyncmodel.ResponseEventConsumer) *view {
	view := &view{
		essentials: es,
		model:      model,
	}

	go view.startEventsProcessing(context.TODO())

	return view
}

func (v *view) startEventsProcessing(ctx context.Context) {
	for event := range v.model.ResponseEvents() {
		err := v.handleEvent(ctx, event)
		if err != nil {
			v.essentials.Logger.Err(err).
				Msgf("failed to handle view event: %s: %+v", event, event)
		}
	}
}

func (v *view) handleEvent(ctx context.Context, event asyncmodel.ResponseEvent) error {
	targetClientIDs := event.TargetClientIDs()

	if len(targetClientIDs) < 3 {
		v.essentials.Logger.Debug().Msgf("sending event %s to %s", event, targetClientIDs)
	} else {
		v.essentials.Logger.Debug().Msgf("sending event %s to %d clients", event, len(targetClientIDs))
	}

	return v.essentials.ClientsHub.SendToClients(
		ctx,
		entities.NewSocketEventPayload(event),
		targetClientIDs...,
	)
}
