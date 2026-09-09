package math

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, Abs(-5))
	assert.Equal(t, 5, Abs(5))
	assert.Equal(t, 10, Abs(10))
}

func TestMinNotZero(t *testing.T) {
	assert.Equal(t, 5, MinNotZero(0, 5))
	assert.Equal(t, 5, MinNotZero(5, 0))
	assert.Equal(t, 3, MinNotZero(3, 8))
}

func TestClamp(t *testing.T) {
	assert.Equal(t, 5, Clamp(2, 5, 10))
	assert.Equal(t, 10, Clamp(15, 5, 10))
	assert.Equal(t, 7, Clamp(7, 5, 10))
}

func TestSubClampZero(t *testing.T) {
	assert.Equal(t, 0, SubClampZero(5, 10))
	assert.Equal(t, 5, SubClampZero(10, 5))
}

func TestDigits(t *testing.T) {
	assert.Equal(t, 1, Digits(0))
	assert.Equal(t, 1, Digits(7))
	assert.Equal(t, 3, Digits(100))
	assert.Equal(t, 4, Digits(-1234))
}

func TestSum(t *testing.T) {
	assert.Equal(t, 15, Sum([]int{1, 2, 3, 4, 5}))
	assert.Equal(t, 0, Sum([]int{}))
}

func TestMaxMap(t *testing.T) {
	t.Run("positive values", func(t *testing.T) {
		m := map[string]int{"a": 10, "b": 50, "c": 30}
		k, v, ok := MaxMap(m)
		assert.True(t, ok)
		assert.Equal(t, "b", k)
		assert.Equal(t, 50, v)
	})

	t.Run("negative values", func(t *testing.T) {
		m := map[string]int{"a": -50, "b": -10, "c": -30}
		k, v, ok := MaxMap(m)
		assert.True(t, ok)
		assert.Equal(t, "b", k)
		assert.Equal(t, -10, v)
	})

	t.Run("empty map", func(t *testing.T) {
		m := map[string]int{}
		_, _, ok := MaxMap(m)
		assert.False(t, ok)
	})
}