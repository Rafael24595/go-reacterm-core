package layout

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestWordNew(t *testing.T) {
	w := New(0, 10)

	assert.Equal(t, 0, w.Start())
	assert.Equal(t, 10, w.End())
}
