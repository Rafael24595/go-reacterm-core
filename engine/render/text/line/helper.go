package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func MaxMeasure(cols winsize.Cols, lines ...Line) winsize.Cols {
	size := winsize.Cols(0)
	for _, l := range lines {
		measure := frag.Measure(cols, l.text...)
		size = max(size, measure)
	}
	return size
}

func HasAtom(atm atom.Atom, lines ...Line) bool {
	for _, line := range lines {
		if frag.HasAtom(atm, line.text...) {
			return true
		}
	}
	return false
}
