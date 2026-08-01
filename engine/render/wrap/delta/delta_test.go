package delta

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

func wordfrag(text string) layout.Frag {
	frg := frag.FromString(text)
	return *layout.NewFrag(&frg)
}

func TestNewDelta(t *testing.T) {
	d := New()

	assert.NotNil(t, d.Frags)
	assert.Empty(t, d.Frags)

	assert.NotNil(t, d.Bounds)
	assert.Empty(t, d.Bounds)

	assert.False(t, d.LeftEdge)
	assert.False(t, d.RightEdge)
}

func TestDelta_Size_AddFrag_BoundAtEnd(t *testing.T) {
	d := New()

	assert.Equal(t, 0, d.Size())

	d.AddFrag(wordfrag("hello"))
	d.AddFrag(wordfrag("world"))

	assert.Equal(t, 2, d.Size())

	d.BoundAtEnd()

	assert.DeepEqual(t, []uint32{2}, d.Bounds)
}

func TestDelta_Merge(t *testing.T) {
	t.Run("merge empty other into non-empty self does nothing", func(t *testing.T) {
		f1 := New()
		f1.AddFrag(wordfrag("a"))
		f1.RightEdge = true

		f2 := New()

		f1.Merge(f2)

		assert.Equal(t, 1, f1.Size())
		assert.True(t, f1.RightEdge)
	})

	t.Run("merge non-empty other into empty self copies state directly", func(t *testing.T) {
		f1 := New()

		f2 := New()
		f2.AddFrag(wordfrag("x"))
		f2.BoundAtEnd()
		f2.LeftEdge = true
		f2.RightEdge = true

		f1.Merge(f2)

		assert.Equal(t, 1, f1.Size())
		assert.DeepEqual(t, []uint32{1}, f1.Bounds)
		assert.True(t, f1.LeftEdge)
		assert.True(t, f1.RightEdge)
	})

	t.Run("merge with matching space boundary (f.RightEdge == other.LeftEdge)", func(t *testing.T) {
		f1 := New()
		f1.AddFrag(wordfrag("a"))
		f1.AddFrag(wordfrag(" "))
		f1.RightEdge = true

		f2 := New()
		f2.AddFrag(wordfrag("b"))
		f2.BoundAtEnd()
		f2.LeftEdge = true
		f2.RightEdge = false

		f1.Merge(f2)

		assert.Equal(t, 3, f1.Size())
		assert.DeepEqual(t, []uint32{3}, f1.Bounds)
		assert.False(t, f1.LeftEdge)
	})

	t.Run("merge with boundary mismatch inserts extra offset bound", func(t *testing.T) {
		f1 := New()
		f1.AddFrag(wordfrag("hello"))
		f1.AddFrag(wordfrag("world"))
		f1.RightEdge = false

		f2 := New()
		f2.AddFrag(wordfrag("foo"))
		f2.BoundAtEnd()
		f2.LeftEdge = true
		f2.RightEdge = true

		f1.Merge(f2)

		assert.Equal(t, 3, f1.Size())
		assert.DeepEqual(t, []uint32{2, 3}, f1.Bounds)
		assert.True(t, f1.RightEdge)
	})
}
