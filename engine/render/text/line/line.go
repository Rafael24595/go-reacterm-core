package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type Line struct {
	Order uint16
	Spec  spec.Spec
	Text  []frag.Frag
	hash  uint64
}

func New(text string, styles ...spec.Spec) *Line {
	return &Line{
		Text: []frag.Frag{
			frag.FromString(text),
		},
		Spec: spec.Merge(styles...),
	}
}

func newLine(
	order uint16,
	spec spec.Spec,
	text []frag.Frag,
) *Line {
	hash := calcHash(
		hash.New(),
		order,
		spec,
		text,
	)

	return &Line{
		Order: order,
		Text:  text,
		Spec:  spec,
		hash:  hash.Sum64(),
	}
}

func calcHash(
	hasher hash.Hasher,
	order uint16,
	spec spec.Spec,
	text []frag.Frag,
) hash.Hasher {
	hasher = hasher.Uint16(order)
	hasher = hasher.Uint64(spec.Hash())
	for _, t := range text {
		hasher = hasher.Uint64(t.Hash())
	}
	return hasher
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

func (l *Line) Clone() *Line {
	spec := l.Spec.Clone()

	text := make([]frag.Frag, len(l.Text))
	copy(text, l.Text)

	return &Line{
		Order: l.Order,
		Spec:  spec,
		Text:  text,
		hash:  l.hash,
	}
}

func Measure(line *Line, cols winsize.Cols) winsize.Cols {
	return spec.Measure(line.Spec, spec.LayoutContext{
		SizeCols: cols,
		TextSize: frag.Measure(cols, line.Text...),
	})
}

func FragsMeasure(cols winsize.Cols, line *Line) winsize.Cols {
	return frag.Measure(cols, line.Text...)
}
