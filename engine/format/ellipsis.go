package format

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Ellipsis struct {
	Data  string
	Count winsize.Cols
}

func NewEllipsis(data string, count winsize.Cols) Ellipsis {
	return Ellipsis{
		Data:  data,
		Count: count,
	}
}

func (e Ellipsis) measure() winsize.Cols {
	return runes.Measure(e.Data) * e.Count
}
