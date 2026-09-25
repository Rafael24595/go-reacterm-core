package format

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestCustomIndex(t *testing.T) {
	t.Run("returns active symbol when cursor equals index", func(t *testing.T) {
		provider := CustomIndex(">", "-")
		gen := provider(0)

		assert.Equal(t, ">", gen(0, 0))
		assert.Equal(t, ">", gen(3, 3))
	})

	t.Run("returns idle symbol when cursor does not match index", func(t *testing.T) {
		provider := CustomIndex(">", "-")
		gen := provider(0)

		assert.Equal(t, "-", gen(0, 1))
		assert.Equal(t, "-", gen(2, 5))
	})
}

func TestPredefinedCustomIndices(t *testing.T) {
	t.Run("GreaterIndex emits active symbol when cursor matches index", func(t *testing.T) {
		gen := GreaterIndex(0)

		assert.Equal(t, "-", gen(0, 0))
		assert.Equal(t, ">", gen(0, 1))
	})

	t.Run("HyphenIndex emits selection pointer when cursor matches index", func(t *testing.T) {
		gen := HyphenIndex(0)

		assert.Equal(t, ">", gen(0, 0))
		assert.Equal(t, "-", gen(0, 1))
	})
}

func TestNumericIndex(t *testing.T) {
	t.Run("NumericListIndex applies default suffix", func(t *testing.T) {
		provider := NumericListIndex()
		gen := provider(2)

		assert.Equal(t, "1 .- ", gen(0, 0))
		assert.Equal(t, "10.- ", gen(0, 9))
	})

	t.Run("NumericListIndex applies default suffix with hight digits value", func(t *testing.T) {
		provider := NumericListIndex()
		gen := provider(3)

		assert.Equal(t, "1  .- ", gen(0, 0))
		assert.Equal(t, "10 .- ", gen(0, 9))
		assert.Equal(t, "100.- ", gen(0, 99))
	})

	t.Run("NumericIndex applies custom suffix", func(t *testing.T) {
		provider := NumericIndex(") ")
		gen := provider(1)

		assert.Equal(t, "1) ", gen(0, 0))
		assert.Equal(t, "5) ", gen(0, 4))
	})

	t.Run("NumericIndex without suffix", func(t *testing.T) {
		provider := NumericIndex()
		gen := provider(2)

		assert.Equal(t, "1 ", gen(0, 0))
	})
}

func TestAlphabeticIndex(t *testing.T) {
	t.Run("AlphabeticListIndex formats sequential letters with default suffix", func(t *testing.T) {
		provider := AlphabeticListIndex()
		gen := provider(1)

		assert.Equal(t, "a.- ", gen(0, 0))
		assert.Equal(t, "b.- ", gen(0, 1))
		assert.Equal(t, "c.- ", gen(0, 2))
	})

	t.Run("AlphabeticListIndex formats sequential letters with default suffix and hight digits value", func(t *testing.T) {
		provider := AlphabeticListIndex()
		gen := provider(2)

		assert.Equal(t, "a .- ", gen(0, 0))
		assert.Equal(t, "b .- ", gen(0, 1))

		assert.Equal(t, "aa.- ", gen(0, 26))
	})

	t.Run("AlphabeticIndex applies custom suffix", func(t *testing.T) {
		provider := AlphabeticIndex(": ")
		gen := provider(1)

		assert.Equal(t, "a: ", gen(0, 0))
		assert.Equal(t, "b: ", gen(0, 1))
	})
}
