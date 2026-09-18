package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type MatrixCursor struct {
	row uint16
	col uint16
	show bool
}

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

func (c *MatrixCursor) IsVisible() bool {
	return c.show
}

func (c *MatrixCursor) WithVisibility(show bool) *MatrixCursor {
	c.show = show
	return c
}

func (c *MatrixCursor) Show() *MatrixCursor {
	c.show = true
	return c
}

func (c *MatrixCursor) Hide() *MatrixCursor {
	c.show = false
	return c
}

func (c *MatrixCursor) Row() uint16 {
	return c.row
}

func (c *MatrixCursor) SetRow(row uint16) *MatrixCursor {
	c.row = row
	return c
}

func (c *MatrixCursor) IncRow(limit uint16) *MatrixCursor {
	c.row = math.AddClampLimit(c.row, 1, limit)
	return c
}

func (c *MatrixCursor) DecRow() *MatrixCursor {
	c.row = math.SubClampZero(c.row, 1)
	return c
}

func (c *MatrixCursor) Col() uint16 {
	return c.col
}

func (c *MatrixCursor) SetCol(col uint16) *MatrixCursor {
	c.col = col
	return c
}

func (c *MatrixCursor) IncCol(limit uint16) *MatrixCursor {
	c.col = math.AddClampLimit(c.col, 1, limit)
	return c
}

func (c *MatrixCursor) DecCol() *MatrixCursor {
	c.col = math.SubClampZero(c.col, 1)
	return c
}

func (c *MatrixCursor) IsAt(row, col uint16) bool {
	if !c.show {
		return false
	}
	return c.row == row && c.col == col
}
