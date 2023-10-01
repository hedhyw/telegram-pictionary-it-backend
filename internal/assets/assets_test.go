package assets_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/assets"
)

func TestWords(t *testing.T) {
	words := assets.Words()

	assert.Contains(t, words, "Flower")
	assert.NotContains(t, words, "")
}

func TestHello(t *testing.T) {
	hello := assets.Hello()

	assert.Contains(t, hello, "Pictionary It")
}
