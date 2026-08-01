package splitter

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

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

func wordToString(word layout.Word, frags []layout.Frag) string {
	return fragsToString(
		frags[word.Start():word.End()],
	)
}

func wordsToStrings(tokens []layout.Word, frags []layout.Frag) []string {
	out := make([]string, len(tokens))
	for i, word := range tokens {
		out[i] = wordToString(word, frags)
	}
	return out
}

func fragsToString(frags []layout.Frag) string {
	var b strings.Builder
	for _, f := range frags {
		b.WriteString(f.Base.Text())
	}
	return b.String()
}
