package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type measureResolver func(winsize.Cols, ...Frag) winsize.Cols

type Frag struct {
	Base     *frag.Frag
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func NewFrag(frg *frag.Frag) *Frag {
	return &Frag{
		Base:     frg,
		measured: false,
		cols:     0,
		measure:  0,
	}
}

func FromFrags(frags ...frag.Frag) []Frag {
	result := make([]Frag, len(frags))
	for i, f := range frags {
		result[i] = *NewFrag(&f)
	}
	return result
}

func (w *Frag) Measure(cols winsize.Cols) winsize.Cols {
	return w.measureWith(cols, fragMeasure)
}

func (w *Frag) measureWith(
	cols winsize.Cols,
	resolver measureResolver,
) winsize.Cols {
	if !w.measured || w.cols != cols {
		w.measure = resolver(cols, *w)
		w.cols = cols
		w.measured = true
	}

	return w.measure
}

func fragMeasure(cols winsize.Cols, frags ...Frag) winsize.Cols {
	measure := winsize.Cols(0)
	for _, f := range frags {
		measure += frag.Measure(cols, *f.Base)
	}
	return measure
}

func AppendFrags(dst []frag.Frag, src []Frag) []frag.Frag {
	for _, f := range src {
		dst = append(dst, *f.Base)
	}
	return dst
}

func splitFragAt(frg *Frag, cols winsize.Cols) (*Frag, *Frag) {
	text := frg.Base.Text()

	if cols == 0 {
		left := clone(frg, "")
		right := clone(frg, text)

		return left, right
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(text, offset.Offset(cols))
	if !canBreak || int(byteIndex) >= len(text) {
		return clone(frg, text), nil
	}

	left := clone(frg, text[:byteIndex])
	right := clone(frg, text[byteIndex:])

	return left, right
}

func clone(frg *Frag, text string) *Frag {
	result := frag.NewBuilder().
		AddText(text).
		WithMeta(*frg.Base).
		Frag()

	return NewFrag(&result)
}
