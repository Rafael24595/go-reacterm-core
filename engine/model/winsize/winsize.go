package winsize

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

// Transformer defines a functional modifier for terminal dimensions.
type Transformer func(Winsize) Winsize

// Rows represents the vertical dimension (height) in lines.
type Rows uint16

// Sub subtracts another Rows value, clamping the result to zero if negative.
func (r Rows) Sub(o Rows) Rows {
	return math.SubClampZero(r, o)
}

// Cols represents the horizontal dimension (width) in characters/cells.
type Cols uint16

// Sub subtracts another Cols value, clamping the result to zero if negative.
func (c Cols) Sub(o Cols) Cols {
	return math.SubClampZero(c, o)
}

// Winsize holds terminal dimensions in rows and columns.
type Winsize struct {
	Rows Rows
	Cols Cols
}

// New constructs a Winsize with the specified rows and columns.
func New(rows Rows, cols Cols) Winsize {
	return Winsize{
		Rows: rows,
		Cols: cols,
	}
}

// Eq checks if two Winsize dimensions are equal.
func (w Winsize) Eq(other Winsize) bool {
	return w == other
}
