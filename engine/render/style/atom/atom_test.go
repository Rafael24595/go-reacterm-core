package atom

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestMerge(t *testing.T) {
	t.Run("Merge multiple flags", func(t *testing.T) {
		merged := Merge(Bold, Focus, Wrap)

		assert.True(t, merged.HasAny(Bold))
		assert.True(t, merged.HasAny(Focus))
		assert.True(t, merged.HasAny(Wrap))
		
		assert.False(t, merged.HasAny(Dim))
	})

	t.Run("Merge empty produces None", func(t *testing.T) {
		merged := Merge()
		assert.Equal(t, None, merged)
	})
}

func TestErase(t *testing.T) {
	initial := Merge(Bold, Focus, Wrap)

	result := Erase(initial, Focus)

	assert.True(t, result.HasAny(Bold))
	assert.False(t, result.HasAny(Focus))
	assert.True(t, result.HasAny(Wrap))
}

func TestAtom_HasAny(t *testing.T) {
	styles := Merge(Bold, Select)

	assert.True(t, styles.HasAny(Bold))
	assert.True(t, styles.HasAny(Select, Dim))
	assert.False(t, styles.HasAny(Dim, Focus))
}

func TestAtom_HasAll(t *testing.T) {
	styles := Merge(Bold, Select, Focus)

	assert.True(t, styles.HasAll(Bold, Select))
	assert.True(t, styles.HasAll(Bold, Focus, Select))
	assert.False(t, styles.HasAll(Bold, Dim))
}

func TestAtom_HasNone(t *testing.T) {
	styles := Merge(Bold, Select)

	assert.True(t, styles.HasNone(Dim, Focus))
	assert.False(t, styles.HasNone(Bold, Dim))
}

func TestAtom_String(t *testing.T) {
	t.Run("None string representation", func(t *testing.T) {
		assert.Equal(t, "None", None.String())
	})

	t.Run("Composite flags string representation", func(t *testing.T) {
		styles := Merge(Bold, Focus)
		assert.Equal(t, "Bold|Focus", styles.String())
	})
}