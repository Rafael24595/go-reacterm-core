package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type Line struct {
	Order uint16
	Text  []frag.Frag
	Spec  spec.Spec
}

func NewLine(text string, styles ...spec.Spec) *Line {
	return &Line{
		Text: []frag.Frag{
			*frag.New(text),
		},
		Spec: spec.Merge(styles...),
	}
}

func EmptyLine(size ...int) *Line {
	bufferSize := 0
	if len(size) > 0 {
		bufferSize = size[0]
	}

	return LineFromFrags(
		make([]frag.Frag, 0, bufferSize)...,
	)
}

func LineFromMeta(other *Line, size ...int) *Line {
	return EmptyLine(size...).
		CopyMeta(other)
}

func LineFromFrags(frags ...frag.Frag) *Line {
	return &Line{
		Text: frags,
		Spec: spec.Empty(),
	}
}

func (l *Line) CopyMeta(other *Line) *Line {
	l.Order = other.Order
	l.AddSpec(other.Spec)
	return l
}

func (l *Line) SetOrder(order uint16) *Line {
	l.Order = order
	return l
}

func (l *Line) UnshiftFrags(frags ...frag.Frag) *Line {
	l.Text = append(frags, l.Text...)
	return l
}

func (l *Line) PushFrags(frags ...frag.Frag) *Line {
	l.Text = append(l.Text, frags...)
	return l
}

func (l *Line) AddSpec(styles ...spec.Spec) *Line {
	newSpec := spec.Merge(styles...)
	l.Spec = spec.Merge(l.Spec, newSpec)
	return l
}

func (l *Line) SetSpec(styles ...spec.Spec) *Line {
	l.Spec = spec.Merge(styles...)
	return l
}

func (l *Line) Clone() *Line {
	newLine := EmptyLine().CopyMeta(l)
	newLine.Text = make([]frag.Frag, len(l.Text))
	copy(newLine.Text, l.Text)
	return newLine
}

func LineMeasure(line *Line, cols winsize.Cols) winsize.Cols {
	return spec.Measure(line.Spec, spec.LayoutContext{
		SizeCols: cols,
		TextSize: frag.Measure(cols, line.Text...),
	})
}
