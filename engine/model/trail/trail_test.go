package trail

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTrail_GoTo(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")

	trail := New(3, mockA)

	trail.GoTo(mockB)

	back, ok := trail.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, "A", back.Name)

	_, ok = trail.PeekForward()
	assert.False(t, ok)
}

func TestTrail_Back(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")
	mockC := screen_test.MockByName("C")

	trail := New(3, mockA)

	trail.GoTo(mockB)
	trail.GoTo(mockC)

	current, ok := trail.Back()
	assert.True(t, ok)
	assert.Equal(t, "B", current.Name)

	forward, ok := trail.PeekForward()
	assert.True(t, ok)
	assert.Equal(t, "C", forward.Name)
}

func TestTrail_Forward(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")
	mockC := screen_test.MockByName("C")

	trail := New(3, mockA)

	trail.GoTo(mockB)
	trail.GoTo(mockC)

	trail.Back()

	current, ok := trail.Forward()
	assert.True(t, ok)
	assert.Equal(t, "C", current.Name)
}

func TestTrail_GoToClearsForward(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")
	mockC := screen_test.MockByName("C")
	mockD := screen_test.MockByName("D")

	trail := New(3, mockA)

	trail.GoTo(mockB)
	trail.GoTo(mockC)

	trail.Back()

	trail.GoTo(mockD)

	_, ok := trail.PeekForward()
	assert.False(t, ok)
}

func TestTrail_RespectsLimit(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")
	mockC := screen_test.MockByName("C")
	mockD := screen_test.MockByName("D")
	mockE := screen_test.MockByName("E")

	trail := New(3, mockA)

	trail.GoTo(mockB)
	trail.GoTo(mockC)
	trail.GoTo(mockD)
	trail.GoTo(mockE)

	current, _ := trail.Back()
	assert.Equal(t, "D", current.Name)

	current, _ = trail.Back()
	assert.Equal(t, "C", current.Name)

	current, _ = trail.Back()
	assert.Equal(t, "B", current.Name)

	_, ok := trail.Back()
	assert.False(t, ok)
}

func TestTrail_PeekDoesNotConsume(t *testing.T) {
	mockA := screen_test.MockByName("A")
	mockB := screen_test.MockByName("B")

	trail := New(3, mockA)

	trail.GoTo(mockB)

	first, ok := trail.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, "A", first.Name)

	second, ok := trail.PeekBack()
	assert.True(t, ok)
	assert.Equal(t, "A", second.Name)

	back, _ := trail.Back()
	assert.Equal(t, "A", back.Name)
}
