package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestTableAction_FluentAPI(t *testing.T) {
	executed := 0

	action := DefaultTableAction().
		WithHandler(func(MatrixCursor) {
			executed += 1
		}).
		AsEdit()

	assert.True(t, action.IsNavigable())
	assert.True(t, action.InEditMode())

	action.Exec(MatrixCursor{})
}

func TestTableAction_DisabledNavigation(t *testing.T) {
	executed := false

	action := DefaultTableAction().
		WithHandler(func(MatrixCursor) {
			executed = false
		}).
		DisableNavigation()

	assert.False(t, action.IsNavigable())

	action.Exec(MatrixCursor{})
	assert.False(t, executed)
}
