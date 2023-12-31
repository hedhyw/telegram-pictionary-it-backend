package player

type Model struct {
	Username string `json:"username"`
	ClientID string `json:"clientId"`

	Score      int `json:"score"`
	RoundScore int `json:"roundScore"`

	IsLead           bool `json:"isLead"`
	RoundWordMatched bool `json:"roundWordMatched"`
}

func New(clientID string, username string) *Model {
	return &Model{
		Username: username,
		ClientID: clientID,
	}
}

func (m *Model) SetRoundWordMatched(roundScore int) {
	m.RoundWordMatched = true
	m.IncRoundScore(roundScore + 1)
}

func (m *Model) IncRoundScore(roundScore int) {
	m.RoundScore += roundScore
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
