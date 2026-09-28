package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestKind_BitmaskOperations(t *testing.T) {
	composite := KindJustifyLeft | KindFill

	t.Run("HasAny", func(t *testing.T) {
		assert.True(t, composite.HasAny(KindJustifyLeft))
		assert.True(t, composite.HasAny(KindFill))
		assert.False(t, composite.HasAny(KindTruncateLeft))
	})

	t.Run("HasAll", func(t *testing.T) {
		assert.True(t, composite.HasAll(KindJustifyLeft, KindFill))
		assert.False(t, composite.HasAll(KindJustifyLeft, KindTruncateLeft))
	})

	t.Run("HasNone", func(t *testing.T) {
		assert.True(t, composite.HasNone(KindTruncateLeft, KindExtendLeft))
		assert.False(t, composite.HasNone(KindJustifyLeft))
	})
}