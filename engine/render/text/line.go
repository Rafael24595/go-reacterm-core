package text

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Line struct {
	Order uint16
	Text  []Fragment
	Spec  spec.Spec
}

func NewLine(text string, styles ...spec.Spec) *Line {
	return &Line{
		Text: []Fragment{
			*NewFragment(text),
		},
		Spec: spec.Merge(styles...),
	}
}

func EmptyLine(size ...int) *Line {
	bufferSize := 0
	if len(size) > 0 {
		bufferSize = size[0]
	}

	return LineFromFragments(
		make([]Fragment, 0, bufferSize)...,
	)
}

func LineFromMeta(other *Line, size ...int) *Line {
	return EmptyLine(size...).
		CopyMeta(other)
}

func LineFromFragments(frags ...Fragment) *Line {
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

func (l *Line) UnshiftFragments(frags ...Fragment) *Line {
	l.Text = append(frags, l.Text...)
	return l
}

func (l *Line) PushFragments(frags ...Fragment) *Line {
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

func (l *Line) CutSpec(styles spec.Kind) *Line {
	l.Spec, _ = spec.Erase(l.Spec, styles)
	return l
}

func (l *Line) Clone() *Line {
	newLine := EmptyLine().CopyMeta(l)
	newLine.Text = make([]Fragment, len(l.Text))
	copy(newLine.Text, l.Text)
	return newLine
}

func LineMeasure(line *Line, cols winsize.Cols) winsize.Cols {
	return spec.Measure(line.Spec, spec.LayoutContext{
		SizeCols: cols,
		TextSize: FragmentMeasure(cols, line.Text...),
	})
}

func LineToString(line *Line) string {
	buffer := make([]string, 0)
	for _, v := range line.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, "")
}
