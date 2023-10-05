package game

import (
	"math/rand"
	"strings"
	"unicode/utf8"

	"github.com/samber/lo"
)

func normalizeWord(word string) string {
	return strings.ToLower(strings.TrimSpace(word))
}

func prepareHint(words string, rand *rand.Rand) string {
	var maxHints int
	size := utf8.RuneCountInString(words)

	switch {
	case size <= 3:
		maxHints = 1
	case size <= 6:
		maxHints = 2
	default:
		maxHints = 3
	}

	return prepareHintForMultipleWords(words, maxHints, rand)
}

func prepareHintForMultipleWords(sentence string, maxHintsCount int, rand *rand.Rand) string {
	words := strings.Split(sentence, " ")

	return strings.Join(lo.Map(words, func(word string, _ int) string {
		return prepareHintForSingleWord(word, maxHintsCount, rand)
	}), " ")
}

func prepareHintForSingleWord(word string, maxHintsCount int, rand *rand.Rand) string {
	if word == "" {
		return ""
	}

	wordRunes := []rune(word)

	hintRunes := []rune(strings.TrimSpace(
		strings.Repeat("_", len(wordRunes)),
	))

	for i := 0; i < maxHintsCount; i++ {
		hintIndex := rand.Intn(len(wordRunes))

		hintRunes[hintIndex] = wordRunes[hintIndex]
	}

	return strings.ToLower(string(hintRunes))
}
