package wrap

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

func lineToString(line layout.Line) string {
	var sb strings.Builder

	for idx := range line.Size() {
		text := fragsToString(
			line.FindFrags(uint(idx)),
		)
		sb.WriteString(text)
	}

	return sb.String()
}

func fragsToString(frags []layout.Frag) string {
	var b strings.Builder
	for _, f := range frags {
		b.WriteString(f.Base.Text())
	}
	return b.String()
}
