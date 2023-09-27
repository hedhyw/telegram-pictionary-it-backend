package assets

import (
	_ "embed"
	"strings"
)

//go:embed words.txt
var words string

func Words() []string {
	return strings.Split(words, "\n")
}
