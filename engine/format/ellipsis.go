package format

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Ellipsis represents a truncation indicator symbol and its repetition count.
type Ellipsis struct {
	// Data is the string used to indicate truncation (e.g., "...").
	Data string
	// Count is the number of times the Data string should be repeated.
	Count winsize.Cols
}

// NewEllipsis returns a custom Ellipsis configuration.
func NewEllipsis(data string, count winsize.Cols) Ellipsis {
	return Ellipsis{
		Data:  data,
		Count: count,
	}
}

// String returns the string representation of the Ellipsis, repeating the Data string Count times.
func (e Ellipsis) string() string {
	if e.Count <= 0 {
		return ""
	}
	return strings.Repeat(e.Data, int(e.Count))
}

func (e Ellipsis) measure() winsize.Cols {
	return runes.Measure(e.Data) * e.Count
}
