package cache

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestMemory_ToCache(t *testing.T) {
	cache := NewMemory[any, any]()
	Helper_ToCache(t, cache)
}

func TestMemory_Operations(t *testing.T) {
	t.Run("Put and Get successfully", func(t *testing.T) {
		c := NewMemory[string, int]()

		c.Put("user_1", 42)

		val, found := c.Get("user_1")

		assert.True(t, found)
		assert.Equal(t, 42, val)
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		c := NewMemory[string, string]()

		_, found := c.Get("non_existent")
		assert.False(t, found)
	})

	t.Run("Delete key", func(t *testing.T) {
		c := NewMemory[int, string]()

		c.Put(100, "laptop")
		c.Del(100)

		_, found := c.Get(100)
		assert.False(t, found)
	})

	t.Run("Len reflects item count", func(t *testing.T) {
		c := NewMemory[string, float64]()

		assert.Equal(t, 0, c.Len())

		c.Put("pi", 3.1416)
		c.Put("e", 2.7182)

		assert.Equal(t, 2, c.Len())

		c.Del("pi")

		assert.Equal(t, 1, c.Len())
	})

	t.Run("Cls clears all items", func(t *testing.T) {
		c := NewMemory[string, string]()

		c.Put("k1", "v1")
		c.Put("k2", "v2")

		c.Cls()

		assert.Equal(t, 0, c.Len())

		_, found := c.Get("k1")
		assert.False(t, found)
	})
}
