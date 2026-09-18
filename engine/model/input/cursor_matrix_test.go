package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewMatrixCursor(t *testing.T) {
	cursor := NewMatrixCursor(2, 4, true)

	assert.Equal(t, 2, cursor.row)
	assert.Equal(t, 4, cursor.col)
	assert.True(t, cursor.show)
}

func TestMatrixCursor_SettersAndVisibility(t *testing.T) {
	cursor := NewMatrixCursor(0, 0, false)

	cursor.SetRow(5).SetCol(10).Show()

	assert.Equal(t, 5, cursor.Row())
	assert.Equal(t, 10, cursor.Col())
	assert.True(t, cursor.IsVisible())

	cursor.Hide()
	assert.False(t, cursor.IsVisible())

	cursor.WithVisibility(true)
	assert.True(t, cursor.IsVisible())
}

func TestMatrixCursor_MovementAndClamping(t *testing.T) {
	cursor := NewMatrixCursor(0, 0, true)

	cursor.IncRow(2)
	assert.Equal(t, 1, cursor.row)

	cursor.IncRow(2)
	assert.Equal(t, 2, cursor.row)

	cursor.IncRow(2)
	assert.Equal(t, 2, cursor.row)

	cursor.DecRow()
	assert.Equal(t, 1, cursor.row)

	cursor.DecRow()
	assert.Equal(t, 0, cursor.row)

	cursor.DecRow()
	assert.Equal(t, 0, cursor.row)

	cursor.IncCol(1)
	assert.Equal(t, 1, cursor.col)

	cursor.IncCol(1)
	assert.Equal(t, 1, cursor.col)

	cursor.DecCol()
	assert.Equal(t, 0, cursor.col)
}

func TestMatrixCursor_IsAt(t *testing.T) {
	cursor := NewMatrixCursor(1, 3, true)

	assert.True(t, cursor.IsAt(1, 3))
	assert.False(t, cursor.IsAt(1, 2))

	cursor.show = false
	assert.False(t, cursor.IsAt(1, 3))
}
