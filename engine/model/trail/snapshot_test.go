package trail

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestSnapshot_ToSlice(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")
	mockC := screen_test.MockByName("C")

	snap := Snapshot{
		Previous: []screen.Node{mockA},
		Current:  mockB,
		Next:     []screen.Node{mockC},
	}

	slice := snap.ToSlice()

	assert.Size(t, 3, slice)
	assert.Equal(t, "A", slice[0].Name)
	assert.Equal(t, "B", slice[1].Name)
	assert.Equal(t, "C", slice[2].Name)
}
