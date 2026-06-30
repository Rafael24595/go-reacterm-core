package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestStack_PushPop_LIFO(t *testing.T) {
	s := New[int](3)

	s.Push(1)
	s.Push(2)
	s.Push(3)

	v, ok := s.Pop()
	assert.True(t, ok)
	assert.Equal(t, 3, v)

	v, ok = s.Pop()
	assert.True(t, ok)
	assert.Equal(t, 2, v)

	v, ok = s.Pop()
	assert.True(t, ok)
	assert.Equal(t, 1, v)

	_, ok = s.Pop()
	assert.False(t, ok)
}

func TestStack_ReplaceWhenFull(t *testing.T) {
	s := New[int](3)

	_, r1 := s.Push(1)
	_, r2 := s.Push(2)
	_, r3 := s.Push(3)

	assert.False(t, r1 || r2 || r3)

	discarded, replaced := s.Push(4)
	assert.True(t, replaced)
	assert.Equal(t, 1, discarded)

	v, _ := s.Pop()
	assert.Equal(t, 4, v)

	v, _ = s.Pop()
	assert.Equal(t, 3, v)

	v, _ = s.Pop()
	assert.Equal(t, 2, v)
}

func TestStack_Clear(t *testing.T) {
	s := New[int](3)

	s.Push(1)
	s.Push(2)

	s.Clear()
	assert.Equal(t, 0, s.Len())

	_, ok := s.Pop()
	assert.False(t, ok)

	s.Push(10)

	v, ok := s.Pop()
	assert.True(t, ok)
	assert.Equal(t, 10, v)
}

func TestStack_Len(t *testing.T) {
	s := New[int](3)

	assert.Equal(t, 0, s.Len())

	s.Push(1)

	assert.Equal(t, 1, s.Len())
}

func TestStack_Capacity(t *testing.T) {
	s := New[int](3)

	assert.Equal(t, 3, s.Cap())

	s = New[int](6)

	assert.Equal(t, 6, s.Cap())
}

func TestStack_Items(t *testing.T) {
	s := New[int](3)

	s.Push(1)
	s.Push(2)
	s.Push(3)

	items := s.Items()

	assert.DeepEqual(t, []int{3, 2, 1}, items)
}

func TestStack_ItemsReturnsCopy(t *testing.T) {
	s := New[int](3)

	s.Push(1)
	s.Push(2)
	s.Push(3)

	items := s.Items()

	items[0] = 100

	again := s.Items()

	assert.DeepEqual(t, []int{3, 2, 1}, again)
}
