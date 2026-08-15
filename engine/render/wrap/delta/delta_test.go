package delta

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewDelta(t *testing.T) {
	d := New()

	assert.NotNil(t, d.frags)
	assert.Empty(t, d.frags)

	assert.NotNil(t, d.bounds)
	assert.Empty(t, d.bounds)

	assert.False(t, d.leftEdge)
	assert.False(t, d.rightEdge)
}
