package hint

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestSizeFixed(t *testing.T) {
	s := Fixed(50)

	assert.Equal(t, 50, s.Min(100))
	assert.Equal(t, 30, s.Min(30))
}

func TestSizePercent(t *testing.T) {
	s := Percent(50)

	assert.Equal(t, 50, s.Min(100))
	assert.Equal(t, 25, s.Min(50))
}

func TestSizeMaximize(t *testing.T) {
	s := Maximize[int]()

	assert.Equal(t, 120, s.Min(120))
}

func TestSizeZeroValueSafety(t *testing.T) {
	var s Size[int]

	assert.Panic(t, func() {
		assert.Equal(t, 0, s.Min(100))
	})
}