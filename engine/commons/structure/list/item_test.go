package list

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestItem_Next(t *testing.T) {
	t.Run("returns next element when present", func(t *testing.T) {
		l := New[int]()
		item1 := l.Push(10)
		item2 := l.Push(20)

		next, ok := item1.Next()

		assert.True(t, ok)
		assert.Equal(t, item2, next)
		assert.Equal(t, 20, next.Data)
	})

	t.Run("returns false when item is the last element", func(t *testing.T) {
		l := New[string]()
		l.Push("A")
		item2 := l.Push("B")

		next, ok := item2.Next()

		assert.False(t, ok)
		assert.Nil(t, next)
	})

	t.Run("returns false when item does not belong to any list", func(t *testing.T) {
		item := newItem("isolated")

		next, ok := item.Next()

		assert.False(t, ok)
		assert.Nil(t, next)
	})

	t.Run("returns false after item is deleted from list", func(t *testing.T) {
		l := New[int]()
		item1 := l.Push(1)
		item2 := l.Push(2)

		l.Delete(item1)

		next, ok := item1.Next()

		assert.False(t, ok)
		assert.Nil(t, next)

		_, ok = item2.Next()
		assert.False(t, ok)
	})
}

func TestItem_Prev(t *testing.T) {
	t.Run("returns previous element when present", func(t *testing.T) {
		l := New[int]()
		item1 := l.Push(100)
		item2 := l.Push(200)

		prev, ok := item2.Prev()

		assert.True(t, ok)
		assert.Equal(t, item1, prev)
		assert.Equal(t, 100, prev.Data)
	})

	t.Run("returns false when item is the first element", func(t *testing.T) {
		l := New[string]()
		item1 := l.Push("head")
		l.Push("tail")

		prev, ok := item1.Prev()

		assert.False(t, ok)
		assert.Nil(t, prev)
	})

	t.Run("returns false when item does not belong to any list", func(t *testing.T) {
		item := newItem(42)

		prev, ok := item.Prev()

		assert.False(t, ok)
		assert.Nil(t, prev)
	})

	t.Run("returns false after item is deleted from list", func(t *testing.T) {
		l := New[int]()
		item1 := l.Push(10)
		item2 := l.Push(20)

		l.Delete(item2)

		prev, ok := item2.Prev()

		assert.False(t, ok)
		assert.Nil(t, prev)

		_, ok = item1.Prev()
		assert.False(t, ok)
	})
}

func TestItem_BidirectionalTraversal(t *testing.T) {
	l := New[int]()
	i1 := l.Push(1)
	i2 := l.Push(2)
	i3 := l.Push(3)

	next1, ok1 := i1.Next()
	assert.True(t, ok1)
	assert.Equal(t, i2, next1)

	next2, ok2 := next1.Next()
	assert.True(t, ok2)
	assert.Equal(t, i3, next2)

	prev2, ok3 := i3.Prev()
	assert.True(t, ok3)
	assert.Equal(t, i2, prev2)

	prev1, ok4 := prev2.Prev()
	assert.True(t, ok4)
	assert.Equal(t, i1, prev1)
}
