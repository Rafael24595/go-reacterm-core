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

func TestBuilder_Size_AddFrag_BoundAtEnd(t *testing.T) {
	b := NewBuilder()

	assert.Equal(t, 0, b.Size())

	b.AddFrag(wordfrag("hello"))
	b.AddFrag(wordfrag("world"))

	assert.Equal(t, 2, b.Size())

	b.BoundAtEnd()

	assert.DeepEqual(t, []uint32{2}, b.Bounds)
}

func TestBuilder_WithDelta(t *testing.T) {
	t.Run("add delta empty other into non-empty self does nothing", func(t *testing.T) {
		b := NewBuilder()
		b.AddFrag(wordfrag("a"))
		b.RightEdge = true

		d := New()

		b.WithDelta(d)

		assert.Equal(t, 1, b.Size())
		assert.True(t, b.RightEdge)
	})

	t.Run("add delta non-empty other into empty self copies state directly", func(t *testing.T) {
		b1 := NewBuilder()

		b2 := NewBuilder().
			AddFrag(wordfrag("x")).
			BoundAtEnd()

		b2.LeftEdge = true
		b2.RightEdge = true

		b1.WithDelta(b2.ToDelta())

		assert.Equal(t, 1, b1.Size())
		assert.DeepEqual(t, []uint32{1}, b1.Bounds)
		assert.True(t, b1.LeftEdge)
		assert.True(t, b1.RightEdge)
	})

	t.Run("add delta with matching space boundary (f.RightEdge == other.LeftEdge)", func(t *testing.T) {
		b1 := NewBuilder().
			AddFrag(wordfrag("a")).
			AddFrag(wordfrag(" ")).
			SetRightEdge(true)

		b2 := NewBuilder().
			AddFrag(wordfrag("b")).
			BoundAtEnd().
			SetLeftEdge(true).
			SetRightEdge(false)

		b1.WithDelta(b2.ToDelta())

		assert.Equal(t, 3, b1.Size())
		assert.DeepEqual(t, []uint32{3}, b1.Bounds)
		assert.False(t, b1.LeftEdge)
	})

	t.Run("add delta with boundary mismatch inserts extra offset bound", func(t *testing.T) {
		b1 := NewBuilder().
			AddFrag(wordfrag("hello")).
			AddFrag(wordfrag("world")).
			SetRightEdge(false)

		b2 := NewBuilder().
			AddFrag(wordfrag("foo")).
			BoundAtEnd().
			SetLeftEdge(true).
			SetRightEdge(true)

		b1.WithDelta(b2.ToDelta())

		assert.Equal(t, 3, b1.Size())
		assert.DeepEqual(t, []uint32{2, 3}, b1.Bounds)
		assert.True(t, b1.RightEdge)
	})
}
