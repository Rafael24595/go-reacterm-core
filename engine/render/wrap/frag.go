package wrap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type measureResolver func(winsize.Cols, ...wordFrag) winsize.Cols

type wordFrag struct {
	Base     *frag.Frag
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func newWordFrag(frg *frag.Frag) *wordFrag {
	return &wordFrag{
		Base:     frg,
		measured: false,
		cols:     0,
		measure:  0,
	}
}

func (w *wordFrag) Measure(cols winsize.Cols) winsize.Cols {
	return w.measureWith(cols, fragMeasure)
}

func (w *wordFrag) measureWith(
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

func toWordFrag(frags ...frag.Frag) []wordFrag {
	result := make([]wordFrag, len(frags))
	for i, f := range frags {
		result[i] = *newWordFrag(&f)
	}
	return result
}

func appendFrags(dst []frag.Frag, src []wordFrag) []frag.Frag {
	for _, f := range src {
		dst = append(dst, *f.Base)
	}
	return dst
}

func fragMeasure(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
	measure := winsize.Cols(0)
	for _, f := range frags {
		measure += frag.Measure(cols, *f.Base)
	}
	return measure
}
