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

type Standard struct {
	atom styler.Atom
	spec styler.Spec
}

func New(atom styler.Atom, spec styler.Spec) Standard {
	return Standard{
		atom: atom,
		spec: spec,
	}
}

func (r Standard) Render(lines []line.Line, size winsize.Winsize) []string {
	buffer := make([]string, len(lines))

	for i, line := range lines {
		text := format.NewText(
			r.renderLineFrags(line, size),
			frag.Measure(size.Cols, line.Text...),
		)

		buffer[i] = r.spec.Apply(line.Spec, size, text)
	}

	return buffer
}

func (r Standard) renderLineFrags(line line.Line, size winsize.Winsize) string {
	var buffer strings.Builder

	frags := ""
	atoms := atom.None

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range line.Text {
		txt := format.NewText(
			f.Text,
			f.Size(),
		)

		spec := r.spec.Apply(f.Spec, lineSize, txt)

		fragSize := frag.Measure(size.Cols, f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		if atoms != f.Atom && len(frags) != 0 {
			atom := r.atom.Apply(frags, atoms)
			buffer.WriteString(atom)

			frags = spec
			atoms = f.Atom

			continue
		}

		frags += spec
		atoms = f.Atom
	}

	if len(frags) != 0 {
		atom := r.atom.Apply(frags, atoms)
		buffer.WriteString(atom)
	}

	return buffer.String()
}
