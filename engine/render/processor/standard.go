package processor

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type standard struct {
	atom styler.Atom
	spec styler.Spec
}

// New constructs a standard processor with the given atom and spec styling rules.
func New(atom styler.Atom, spec styler.Spec) standard {
	return standard{
		atom: atom,
		spec: spec,
	}
}

// Render processes all lines and returns a single newline-delimited string ready for output.
func (r standard) Render(lines []line.Line, size winsize.Winsize) string {
	buffer := r.RawRender(lines, size)
	return strings.Join(buffer, "\n")
}

// RawRender converts each line into its corresponding styled string representation.
func (r standard) RawRender(lines []line.Line, size winsize.Winsize) []string {
	if len(lines) == 0 {
		return []string{}
	}

	buffer := make([]string, len(lines))

	for i, lne := range lines {
		text := format.NewText(
			r.renderLineFrags(lne, size),
			line.FragsMeasure(size.Cols, lne),
		)

		buffer[i] = r.spec.Apply(lne.Spec(), size, text)
	}

	return buffer
}

func (r standard) renderLineFrags(line line.Line, size winsize.Winsize) string {
	var lineBuffer strings.Builder
	var fragBuffer strings.Builder

	currentAtom := atom.None

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for f := range line.All() {
		txt := format.NewText(
			f.Text(),
			f.Measure(),
		)

		spec := r.spec.Apply(f.Spec(), lineSize, txt)

		fragSize := frag.Measure(size.Cols, *f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		if currentAtom != f.Atom() && fragBuffer.Len() != 0 {
			atom := r.atom.Apply(fragBuffer.String(), currentAtom)
			lineBuffer.WriteString(atom)

			fragBuffer.Reset()
		}

		fragBuffer.WriteString(spec)
		currentAtom = f.Atom()
	}

	if fragBuffer.Len() != 0 {
		atom := r.atom.Apply(fragBuffer.String(), currentAtom)
		lineBuffer.WriteString(atom)
	}

	return lineBuffer.String()
}
