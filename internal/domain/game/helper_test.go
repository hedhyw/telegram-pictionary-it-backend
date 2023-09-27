package game

import (
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
