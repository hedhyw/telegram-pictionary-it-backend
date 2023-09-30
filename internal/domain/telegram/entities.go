package telegram

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

type InitDataMeta struct {
	UserJSON     *string `schema:"user"`
	ChatInstance string  `schema:"chat_instance"`
	ChatType     string  `schema:"chat_type"`

	Raw url.Values `schema:"-"`
}

// User unmarshalls UserJSON to WebAppUser if possible.
func (m InitDataMeta) User() (*WebAppUser, error) {
	if m.UserJSON == nil {
		return nil, semerr.NewNotFoundError(semerr.Error("user is nil"))
	}

	var user WebAppUser

	err := json.Unmarshal([]byte(*m.UserJSON), &user)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	return &user, nil
}

type WebAppUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
