package game

import (
	"context"
	"runtime/debug"

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

	view.startEventsProcessing(context.Background())

	return view
}

func (v *view) startEventsProcessing(ctx context.Context) {
	for i := 0; i < v.essentials.Config.WorkersCount; i++ {
		go func(i int) {
			logger := v.essentials.Logger.With().Int("worker", i).Logger()

			for event := range v.model.ResponseEvents() {
				err := v.handleEvent(ctx, event)
				if err != nil {
					logger.Err(err).
						Msgf("failed to handle view event: %s: %+v", event, event)
				}
			}
		}(i)
	}
}

func (v *view) handleEvent(ctx context.Context, event asyncmodel.ResponseEvent) error {
	defer func() {
		if r := recover(); r != nil {
			v.essentials.Logger.Error().Msgf("panic %v", r)
			debug.PrintStack()
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, v.essentials.Config.ServerTimeout)
	defer cancel()

	targetClientIDs := event.TargetClientIDs()

	if len(targetClientIDs) < 3 {
		v.essentials.Logger.Debug().Msgf("sending event %s to %s", event, targetClientIDs)
	} else {
		v.essentials.Logger.Debug().Msgf("sending event %s to %d clients", event, len(targetClientIDs))
	}

	return v.essentials.ClientsHub.SendToClients(
		ctx,
		entities.NewSocketResponse(event),
		targetClientIDs...,
	)
}
