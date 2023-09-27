package player

import (
	"github.com/brianvoe/gofakeit"
)

type Model struct {
	Username string `json:"username"`
	ClientID string `json:"clientId"`

	Score      int `json:"score"`
	RoundScore int `json:"roundScore"`

	IsLead           bool `json:"isLead"`
	RoundWordMatched bool `json:"roundWordMatched"`
}

func New(clientID string) *Model {
	return &Model{
		// TODO: check username unique.
		Username: gofakeit.Username(),
		ClientID: clientID,
	}
}

func (m *Model) SetRoundWordMatched() {
	m.RoundWordMatched = true
	m.RoundScore++
	m.Score += m.RoundScore
}

func (m *Model) SetLeader() {
	m.IsLead = true
}

func (m *Model) ResetRound() {
	m.IsLead = false
	m.RoundScore = 0
	m.RoundWordMatched = false
}
