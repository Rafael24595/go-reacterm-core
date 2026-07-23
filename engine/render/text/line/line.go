package line

import (
	"iter"

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

func New(
	order uint16,
	spec spec.Spec,
	text []frag.Frag,
) Line {
	hash := calcHash(
		hash.New(),
		order,
		spec,
		text,
	)

	return Line{
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

func (l *Line) Size() uint {
	return uint(len(l.Text))
}

func (l *Line) GetOrder() uint16 {
	return l.Order
}

func (l *Line) GetSpec() spec.Spec {
	return l.Spec
}

func (l *Line) GetText() []frag.Frag {
	return l.Text
}

func (l *Line) GetFrag(index uint) frag.Frag {
	frg, _ := l.TryGetFrag(index)
	return frg
}

func (l *Line) TryGetFrag(index uint) (frag.Frag, bool) {
	if index >= l.Size() {
		return frag.Frag{}, false
	}
	return l.Text[index], true
}

func (l *Line) Frags() iter.Seq[frag.Frag] {
	return func(yield func(frag.Frag) bool) {
		for _, f := range l.Text {
			if !yield(f) {
				return
			}
		}
	}
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

func FragsMeasure(cols winsize.Cols, line Line) winsize.Cols {
	return frag.Measure(cols, line.Text...)
}
