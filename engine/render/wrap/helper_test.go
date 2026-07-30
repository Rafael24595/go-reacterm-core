package wrap

import (
	"strings"
)

func wordsToString(words []word, frags []wordFrag) string {
	var sb strings.Builder
	for _, words := range words {
		sb.WriteString(
			wordToString(words, frags),
		)
	}
	return sb.String()
}

func wordToString(word word, frags []wordFrag) string {
	return fragsToString(frags[word.start:word.end])
}

func wordsToStrings(tokens []word, frags []wordFrag) []string {
	out := make([]string, len(tokens))
	for i, word := range tokens {
		out[i] = wordToString(word, frags)
	}
	return out
}

func fragsToString(frags []wordFrag) string {
	var b strings.Builder
	for _, f := range frags {
		b.WriteString(f.Base.Text())
	}
	return b.String()
}
