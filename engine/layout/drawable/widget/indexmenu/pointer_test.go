package indexmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPointer_Navigation(t *testing.T) {
	t.Run("findPointer bounds checking", func(t *testing.T) {
		assert.Equal(t, pointerSelect, FindPointer(0))
		assert.Equal(t, pointerBold, FindPointer(1))

		assert.Equal(t, pointerSelect, FindPointer(3))
		assert.Equal(t, pointerSelect, FindPointer(255))
	})

	t.Run("nextPointer cycling logic", func(t *testing.T) {
		assert.Equal(t, uint8(1), NextPointer(0))
		assert.Equal(t, uint8(0), NextPointer(1))
		assert.Equal(t, uint8(1), NextPointer(2))
	})
}
