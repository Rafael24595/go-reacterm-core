package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Word struct {
	start    uint32
	end      uint32
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func New(start uint32, end uint32) *Word {
	return &Word{
		start: start,
		end:   end,
	}
}

func (w *Word) Start() uint32 {
	return w.start
}

func (w *Word) End() uint32 {
	return w.end
}
