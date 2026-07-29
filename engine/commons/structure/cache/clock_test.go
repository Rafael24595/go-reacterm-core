package cache

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestClock_ToCache(t *testing.T) {
	cache := NewClock[any, any]()
	Helper_ToCache(t, cache)
}

func TestClock_Initialization(t *testing.T) {
	t.Run("default capacity", func(t *testing.T) {
		c := NewClock[string, int]()

		assert.Equal(t, 0, c.Len())
	})

	t.Run("custom capacity", func(t *testing.T) {
		c := NewClock[string, int](10)
		assert.Equal(t, 0, c.Len())
	})
}

func TestClock_PutAndGet(t *testing.T) {
	t.Run("insert and retrieve items", func(t *testing.T) {
		c := NewClock[string, string](5)

		c.Put("a", "alpha")
		c.Put("b", "beta")

		val, found := c.Get("a")

		assert.True(t, found)
		assert.Equal(t, "alpha", val)

		val, found = c.Get("b")

		assert.True(t, found)
		assert.Equal(t, "beta", val)
	})

	t.Run("update existing key", func(t *testing.T) {
		c := NewClock[string, int](3)

		c.Put("k1", 10)
		c.Put("k1", 20)

		assert.Equal(t, 1, c.Len())

		val, found := c.Get("k1")

		assert.True(t, found)
		assert.Equal(t, 20, val)
	})

	t.Run("get non-existent key returns zero value", func(t *testing.T) {
		c := NewClock[string, int](3)

		_, found := c.Get("missing")
		assert.False(t, found)
	})
}

func TestClock_Del(t *testing.T) {
	t.Run("delete existing key", func(t *testing.T) {
		c := NewClock[string, string](5)

		c.Put("k1", "v1")
		c.Put("k2", "v2")

		c.Del("k1")

		assert.Equal(t, 1, c.Len())

		_, found := c.Get("k1")
		assert.False(t, found)

		val, found := c.Get("k2")
		assert.True(t, found)
		assert.Equal(t, "v2", val)
	})

	t.Run("delete non-existing key does nothing", func(t *testing.T) {
		c := NewClock[string, int](3)

		c.Put("a", 100)
		c.Del("non_existent")

		assert.Equal(t, 1, c.Len())
	})

	t.Run("put and evict work correctly after deletion (no nil pointer panic)", func(t *testing.T) {
		c := NewClock[string, string](3)

		c.Put("k1", "v1")
		c.Put("k2", "v2")
		c.Put("k3", "v3")

		c.Del("k2")

		c.Put("k4", "v4")
		c.Put("k5", "v5")

		assert.Equal(t, 3, c.Len())
	})

	t.Run("delete all items one by one", func(t *testing.T) {
		c := NewClock[int, string](3)

		c.Put(1, "one")
		c.Put(2, "two")

		c.Del(1)
		c.Del(2)

		assert.Equal(t, 0, c.Len())
	})
}

func TestClock_Eviction(t *testing.T) {
	t.Run("evicts oldest unreferenced item when full", func(t *testing.T) {
		c := NewClock[string, string](3)

		c.Put("k1", "v1")
		c.Put("k2", "v2")
		c.Put("k3", "v3")

		c.Put("k4", "v4")

		assert.Equal(t, 3, c.Len())

		_, found := c.Get("k1")
		assert.False(t, found)

		for _, key := range []string{"k2", "k3", "k4"} {
			_, found := c.Get(key)
			assert.True(t, found)
		}
	})

	t.Run("second chance mechanism saves frequently accessed items", func(t *testing.T) {
		c := NewClock[string, string](3)

		c.Put("k1", "v1")
		c.Put("k2", "v2")
		c.Put("k3", "v3")

		c.Get("k1")

		c.Put("k4", "v4")

		_, found := c.Get("k1")
		assert.True(t, found)

		_, found = c.Get("k2")
		assert.False(t, found)
	})
}

func TestClock_Cls(t *testing.T) {
	c := NewClock[string, int](3)

	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)

	c.Cls()

	assert.Equal(t, 0, c.Len())

	for _, key := range []string{"a", "b", "c"} {
		_, found := c.Get(key)
		assert.False(t, found)
	}

	c.Put("d", 4)

	val, found := c.Get("d")

	assert.True(t, found)
	assert.Equal(t, 4, val)
}

func TestClock_Stress(t *testing.T) {
	capacity := uint(5)
	c := NewClock[int, string](capacity)

	for i := range 50 {
		c.Put(i, fmt.Sprintf("val-%d", i))

		assert.LessOrEqual(t, capacity, c.Len())
	}
}
