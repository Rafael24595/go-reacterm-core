package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestCheckAction_FluentAPI(t *testing.T) {
	executed := 0

	action := DefaultCheckAction().
		AsEdit().
		WithHandler(func() {
			executed += 1
		})

	assert.True(t, action.InEditMode())

	action.Exec()
	assert.Equal(t, 1, executed)
}

func TestCheckAction_Empty(t *testing.T) {
	action := DefaultCheckAction()

	assert.False(t, action.InEditMode())
	action.Exec()
}
