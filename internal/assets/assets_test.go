package assets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/telegram-pictionary-backend/internal/assets"
)

func TestWords(t *testing.T) {
	words := assets.Words()

	assert.Contains(t, words, "Flower")
	assert.NotContains(t, words, "")
}
