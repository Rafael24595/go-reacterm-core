package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

// MatrixCursor represents the position and visibility state of a cursor within a matrix or table layout.
type MatrixCursor struct {
	row uint16
	col uint16
	show bool
}

// NewMatrixCursor initializes and returns a pointer to a new MatrixCursor.
func NewMatrixCursor(
	row, col uint16,
	show bool,
) *MatrixCursor {
	return &MatrixCursor{
		row:  row,
		col:  col,
		show: show,
	}
}

// IsVisible returns true if the cursor is currently visible, otherwise false.
func (c *MatrixCursor) IsVisible() bool {
	return c.show
}

// Showing sets the cursor's visibility state and returns the cursor instance for chaining.
func (c *MatrixCursor) WithVisibility(show bool) *MatrixCursor {
	c.show = show
	return c
}

// Show sets the cursor's visibility to true and returns the cursor instance for chaining.
func (c *MatrixCursor) Show() *MatrixCursor {
	c.show = true
	return c
}

// Hide sets the cursor's visibility to false and returns the cursor instance for chaining.
func (c *MatrixCursor) Hide() *MatrixCursor {
	c.show = false
	return c
}

// Rows returns the current row index of the cursor.
func (c *MatrixCursor) Row() uint16 {
	return c.row
}

// SetRow sets the cursor's row index to the specified value and returns the cursor instance for chaining.
func (c *MatrixCursor) SetRow(row uint16) *MatrixCursor {
	c.row = row
	return c
}

// IncRow increments the row index and clamps it to maxRow limit.
func (c *MatrixCursor) IncRow(limit uint16) *MatrixCursor {
	c.row = math.AddClampLimit(c.row, 1, limit)
	return c
}

// DecRow decrements the row index, clamping at zero to prevent underflow.
func (c *MatrixCursor) DecRow() *MatrixCursor {
	c.row = math.SubClampZero(c.row, 1)
	return c
}

// Cols returns the current column index of the cursor.
func (c *MatrixCursor) Col() uint16 {
	return c.col
}

// SetCol sets the cursor's column index to the specified value and returns the cursor instance for chaining.
func (c *MatrixCursor) SetCol(col uint16) *MatrixCursor {
	c.col = col
	return c
}

// IncCol increments the column index and clamps it to maxCol limit.
func (c *MatrixCursor) IncCol(limit uint16) *MatrixCursor {
	c.col = math.AddClampLimit(c.col, 1, limit)
	return c
}

// DecCol decrements the column index, clamping at zero to prevent underflow.
func (c *MatrixCursor) DecCol() *MatrixCursor {
	c.col = math.SubClampZero(c.col, 1)
	return c
}

// IsAt reports whether the cursor is currently visible and positioned at the target row and column.
func (c *MatrixCursor) IsAt(row, col uint16) bool {
	if !c.show {
		return false
	}
	return c.row == row && c.col == col
}
