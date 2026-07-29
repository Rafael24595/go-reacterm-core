package cache

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNewEntry(t *testing.T) {
	e := newEntry("user_key", "user_value")

	assert.Equal(t, "user_key", e.key)
	assert.Equal(t, "user_value", e.value)
	assert.Equal(t, 1, e.used)
}

func TestEntry_Touch(t *testing.T) {
	t.Run("increments used counter", func(t *testing.T) {
		e := newEntry("k", "v")
		e.touch(5)

		assert.Equal(t, 2, e.used)
	})

	t.Run("respects maximum threshold", func(t *testing.T) {
		e := newEntry("k", "v")

		maxLimit := uint8(3)

		e.touch(maxLimit)
		e.touch(maxLimit)
		e.touch(maxLimit)

		assert.Equal(t, maxLimit, e.used)
	})

	t.Run("prevents overflow", func(t *testing.T) {
		e := newEntry("k", "v")
		e.used = ^uint8(0)

		maxLimit := uint8(3)

		e.touch(maxLimit)

		assert.Equal(t, maxLimit, e.used)
	})
}

func TestEntry_Cool(t *testing.T) {
	t.Run("decrements used counter and returns false when used > 0", func(t *testing.T) {
		e := newEntry("k", "v")

		isCold := e.cool()

		assert.False(t, isCold)
		assert.Equal(t, 0, e.used)
	})

	t.Run("returns true when used reaches 0 and does not underflow", func(t *testing.T) {
		e := newEntry("k", "v")

		e.cool()
		isCold := e.cool()

		assert.True(t, isCold)
		assert.Equal(t, 0, e.used)
	})
}
