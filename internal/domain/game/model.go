package game

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/samber/lo"

	"github.com/hedhyw/telegram-pictionary-backend/internal/assets"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/player"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/telegram"
)

type Model struct {
	essentials Essentials

	*asyncmodel.Model

	leaderIndex int
	players     []*player.Model
	word        string
	round       int
	finishAt    time.Time

	roundDoneCh        chan struct{}
	autoFinisherDoneCh chan struct{}
}

func newModel(es Essentials) *Model {
	model := &Model{
		essentials: es,

		Model: nil,

		leaderIndex: 0,
		players:     []*player.Model{},
		word:        "",
		round:       0,
	}

	model.Model = asyncmodel.New(asyncmodel.Essentials{
		InitialState:           &stateInitial{model: model},
		HandleRequestErrorFunc: asyncmodel.DefaultLogRequestErrorHandler(es.Logger),
		RequestTimeout:         es.Config.ServerTimeout,
		Logger:                 es.Logger,
		ChannelSize:            es.Config.WorkersCount,
	})

	return model
}

func (m *Model) addPlayer(clientID string, meta *telegram.InitDataMeta) *player.Model {
	logger := m.essentials.Logger

	username := makeUsernameUnique(
		getUsername(meta),
		m.players,
	)

	player := player.New(clientID, username)
	m.players = append(m.players, player)

	logger.Debug().
		Str("client", clientID).
		Msgf("the player %s joined the game in the chat %s", player.Username, m.essentials.ChatID)

	return player
}

func makeUsernameUnique(username string, players []*player.Model) string {
	const limitCombinations = 1_000

	usernamesSet := lo.SliceToMap(players, func(player *player.Model) (string, struct{}) {
		return player.Username, struct{}{}
	})

	_, ok := usernamesSet[username]
	if !ok {
		return username
	}

	for i := 2; i < limitCombinations; i++ {
		usernameWithNumber := username + " " + strconv.Itoa(i)

		_, ok := usernamesSet[usernameWithNumber]
		if !ok {
			return usernameWithNumber
		}
	}

	return username
}

func getUsername(meta *telegram.InitDataMeta) string {
	user, err := meta.User()
	if err != nil {
		return gofakeit.Username()
	}

	if user.Username != "" {
		return user.Username
	}

	if user.FirstName != "" || user.LastName != "" {
		return strings.TrimSpace(user.FirstName + " " + user.LastName)
	}

	return gofakeit.Username()
}

func (m *Model) setRandomWord() {
	words := assets.Words()
	//nolint: gosec // It is a game.
	index := rand.Intn(len(words))
	m.word = words[index]
}

func (m *Model) responseEventGameStateChanged() *ResponseEventGameStateChanged {
	return &ResponseEventGameStateChanged{
		Players:  m.players,
		State:    m.State(),
		FinishAt: m.finishAt.UTC(),
	}
}

func (m *Model) isEveryoneGuessed() bool {
	for _, player := range m.players {
		if !player.IsLead && !player.RoundWordMatched {
			return false
		}
	}

	return true
}

func (m *Model) getLeader() *player.Model {
	if len(m.players) == 0 {
		return nil
	}

	return m.players[m.leaderIndex]
}

func (m *Model) finishGame(ctx context.Context) error {
	logger := m.essentials.Logger

	close(m.roundDoneCh)

	select {
	case <-m.autoFinisherDoneCh:
		logger.Debug().Msg("game is finished")
	case <-ctx.Done():
		return ctx.Err()
	}

	m.SetState(&stateFinished{model: m})

	return m.EmitResponses(ctx, m.responseEventGameStateChanged())
}

func (m *Model) runGameAutoFinisher() {
	finishOnce := &sync.Once{}
	defer finishOnce.Do(func() { close(m.autoFinisherDoneCh) })
	logger := m.essentials.Logger

	timer := time.NewTimer(m.essentials.Config.GameRoundTimeout)
	defer timer.Stop()

	select {
	case <-timer.C:
		ctx := context.Background()

		logger.Debug().Msg("game is timeouted")

		finishOnce.Do(func() { close(m.autoFinisherDoneCh) })

		err := m.finishGame(ctx)
		if err != nil {
			logger.Err(err).Msg("failed to finish the game")
		}
	case <-m.roundDoneCh:
		logger.Debug().Msg("game auto finisher stopped")
	}
}

func (m *Model) removePlayer(
	ctx context.Context,
	clientID string,
) error {
	logger := m.essentials.Logger.With().Str("client", clientID).Logger()

	m.players = lo.Filter(m.players, func(item *player.Model, index int) bool {
		return item.ClientID != clientID
	})

	logger.Debug().Msgf("removed player")

	return m.EmitResponses(ctx, m.responseEventGameStateChanged())
}

func (m *Model) startGame(ctx context.Context) error {
	if len(m.players) <= 1 {
		return semerr.NewBadRequestError(errNotEnoughPlayers)
	}

	m.roundDoneCh = make(chan struct{})
	m.autoFinisherDoneCh = make(chan struct{})

	m.finishAt = time.Now().Add(m.essentials.Config.GameRoundTimeout)

	// nolint: contextcheck // Worker has a different context.
	go m.runGameAutoFinisher()

	m.round++

	m.setRandomWord()
	m.leaderIndex = (m.round + 1) % len(m.players)

	for _, p := range m.players {
		p.ResetRound()
	}

	m.players[m.leaderIndex].SetLeader()
	m.SetState(&stateInProgress{model: m})

	return m.EmitResponses(ctx,
		&ResponseEventCanvasChanged{
			Players:       m.players,
			ActorClientID: m.getLeader().ClientID,
			ImageBase64:   "",
		},
		&ResponseEventGameStarted{
			Players: m.players,
		},
		m.responseEventGameStateChanged(),
		&ResponseEventLeadHello{
			ClientID: m.getLeader().ClientID,
			Word:     m.word,
		},
	)
}
