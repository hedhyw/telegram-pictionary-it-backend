package assets

import (
	_ "embed"
	"strings"
)

//go:embed words.txt
var words string

// Words returns all lines from words.txt.
func Words() []string {
	return strings.Split(words, "\n")
}

//go:embed hello.html
var hello string

// Hello returns the content of hello.html.
func Hello() string {
	return hello
}
