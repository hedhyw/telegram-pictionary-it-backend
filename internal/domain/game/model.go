package game

import (
	"context"
	"math/rand"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-backend/internal/assets"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/player"
)

type Model struct {
	essentials Essentials

	*asyncmodel.Model

	leaderIndex int
	players     []*player.Model
	word        string
	round       int
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

	model.Model = asyncmodel.New(
		&stateInitial{model: model},
		asyncmodel.DefaultLogRequestErrorHandler(es.Logger),
	)

	return model
}

func (m *Model) addPlayer(clientID string) *player.Model {
	logger := m.essentials.Logger

	player := player.New(clientID)
	m.players = append(m.players, player)

	logger.Debug().
		Str("client", clientID).
		Msgf("the player %s joined the game in the chat %s", player.Username, m.essentials.ChatID)

	return player
}

func (m *Model) setRandomWord() {
	words := assets.Words()
	//nolint: gosec // It is a game.
	index := rand.Intn(len(words))
	m.word = words[index]
}

func (m *Model) responseEventGameStateChanged() *ResponseEventGameStateChanged {
	return &ResponseEventGameStateChanged{
		Players: m.players,
		State:   m.State(),
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

func (m *Model) startGame(ctx context.Context) error {
	if len(m.players) <= 1 {
		return semerr.NewBadRequestError(errNotEnoughPlayers)
	}

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
