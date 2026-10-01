package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Word represents a tokenized visual text unit tracked by byte boundary indices and cell column measurements.
type Word struct {
	start    uint32
	end      uint32
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

// New constructs a new Word instance with the specified start and end byte offsets.
func New(start uint32, end uint32) *Word {
	return &Word{
		start: start,
		end:   end,
	}
}

// Start returns the starting byte offset of the word within its parent fragment or line.
func (w *Word) Start() uint32 {
	return w.start
}

// End returns the ending byte offset of the word within its parent fragment or line.
func (w *Word) End() uint32 {
	return w.end
}
