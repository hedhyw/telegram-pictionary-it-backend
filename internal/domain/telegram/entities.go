package telegram

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

// InitDataMeta is a structure with input data transferred to the Mini App.
type InitDataMeta struct {
	// UserJSON is an optional JSON string containing data about the current user.
	// Use InitDataMeta.User to decode it to *WebAppUser.
	UserJSON *string `schema:"user"`
	// ChatInstance is a global identifier, uniquely corresponding to
	// the chat from which the Mini App was opened. Returned only for
	// Mini Apps launched from a direct link.
	ChatInstance string `schema:"chat_instance"`
	// ChatType is a type of the chat from which the Mini App was opened. Can be
	// either “sender” for a private chat with the user opening the link,
	// “private”, “group”, “supergroup”, or “channel”. Returned only
	// for Mini Apps launched from direct links.
	ChatType string `schema:"chat_type"`

	// Raw data transferred to the Mini App. It is used for validation.
	Raw url.Values `schema:"-"`
}

// User unmarshalls UserJSON to WebAppUser if it is not nil.
func (m InitDataMeta) User() (*WebAppUser, error) {
	if m.UserJSON == nil {
		return nil, semerr.NewNotFoundError(semerr.Error("user is nil"))
	}

	var user WebAppUser

	err := json.Unmarshal([]byte(*m.UserJSON), &user)
	if err != nil {
		return nil, semerr.NewBadRequestError(fmt.Errorf("unmarshalling: %w", err))
	}

	return &user, nil
}

// WebAppUser is a structure containing the information about Telegram user.
type WebAppUser struct {
	// A unique identifier for the user. It has at most
	// 52 significant bits.
	ID int64 `json:"id"`
	// Username of the user. It is optional.
	Username string `json:"username"`
	// FirstName of the user.
	FirstName string `json:"first_name"`
	// LastName of the user. It is optional.
	LastName string `json:"last_name"`
}
