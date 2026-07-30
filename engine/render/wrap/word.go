package wrap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type word struct {
	start    uint32
	end      uint32
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func newWord(start uint32, end uint32) *word {
	return &word{
		start: start,
		end:   end,
	}
}
