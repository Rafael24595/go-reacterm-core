package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func assembleLines(t *testing.T, lines ...line.Line) string {
	t.Helper()

	var sb strings.Builder

	for i, l := range lines {
		_, err := sb.WriteString(
			text_test.LineToString(l),
		)

		assert.Nil(t, err)

		if i < len(lines)-1 {
			_, err := sb.WriteString("\n")
			assert.Nil(t, err)
		}
	}

	return sb.String()
}

func benchmarkText(size int) string {
	const sample = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. "

	var b strings.Builder
	b.Grow(size)

	for b.Len() < size {
		b.WriteString(sample)
	}

	return b.String()[:size]
}

func benchmarkLine(size int) line.Line {
	return line.FromFrags(
		frag.FromStrings(
			benchmarkText(size),
		)...,
	)
}


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
