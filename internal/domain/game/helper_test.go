package game

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeWord(t *testing.T) {
	t.Parallel()

	testCases := [...]struct {
		Word     string
		Expected string
	}{{
		Word:     "hello",
		Expected: "hello",
	}, {
		Word:     "HELLO",
		Expected: "hello",
	}, {
		Word:     "  Hello  ",
		Expected: "hello",
	}, {
		Word:     "  Hello!  ",
		Expected: "hello!",
	}, {
		Word:     "  Hello ?",
		Expected: "hello ?",
	}}

	for _, testCase := range testCases {
		actual := normalizeWord(testCase.Expected)
		assert.Equal(t, testCase.Expected, actual)
	}
}

func TestPrepareHint(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		word := prepareHint("", newZeroSeedRand())
		assert.Empty(t, word)
	})

	t.Run("short", func(t *testing.T) {
		t.Parallel()

		word := prepareHint("cat", newZeroSeedRand())
		assert.Equal(t, "c__", word)
	})

	t.Run("medium", func(t *testing.T) {
		t.Parallel()

		word := prepareHint("computer", newZeroSeedRand())
		assert.Equal(t, "_om_____", word)
	})

	t.Run("long", func(t *testing.T) {
		t.Parallel()

		word := prepareHint("hello world", newZeroSeedRand())
		assert.Equal(t, "___lo wo___", word)
	})
}

func TestPrepareHintForSinprepareHintForSingleWord(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("", 10, newZeroSeedRand())
		assert.Empty(t, hint)
	})

	t.Run("no_hints", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("abc", 0, newZeroSeedRand())
		assert.Equal(t, "___", hint)
	})

	t.Run("one_hint", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("abc", 1, newZeroSeedRand())
		assert.Equal(t, "a__", hint)
	})

	t.Run("all_hints", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("abc", 10, newZeroSeedRand())
		assert.Equal(t, "abc", hint)
	})

	t.Run("negative_hints_count", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("abc", -1, newZeroSeedRand())
		assert.Equal(t, "___", hint)
	})

	t.Run("lowercase", func(t *testing.T) {
		t.Parallel()

		hint := prepareHintForSingleWord("ABC", 16, newZeroSeedRand())
		assert.Equal(t, "abc", hint)
	})
}

func newZeroSeedRand() *rand.Rand {
	// nolint: gosec // It is a test.
	return rand.New(rand.NewSource(0))
}
