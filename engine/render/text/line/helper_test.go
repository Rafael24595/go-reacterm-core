package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestMaxMeasure(t *testing.T) {
	cols := winsize.Cols(80)

	t.Run("returns 0 when no lines are provided", func(t *testing.T) {
		got := MaxMeasure(cols)
		assert.Equal(t, 0, got)
	})

	t.Run("returns maximum measure among provided lines", func(t *testing.T) {
		l1 := FromString("lne_01")
		l2 := FromFrags(
			frag.TextSpec(
				"frg_02", spec.JustifyCenter(50),
			),
		)

		got := MaxMeasure(cols, l1, l2)

		assert.Equal(t, 50, got)
	})
}

func TestHasAtom(t *testing.T) {
	targetAtom := atom.Bold
	otherAtom := atom.Dim

	t.Run("returns true when at least one line contains the target atom", func(t *testing.T) {
		l1 := FromFrags(
			frag.FromAtom(targetAtom),
		)
		l2 := FromFrags(
			frag.FromAtom(otherAtom),
		)

		assert.True(t, HasAtom(targetAtom, l1, l2))
	})

	t.Run("returns false when no lines contain the target atom", func(t *testing.T) {
		l := FromFrags(
			frag.FromAtom(otherAtom),
		)

		assert.False(t, HasAtom(targetAtom, l))
	})

	t.Run("returns false when line list is empty", func(t *testing.T) {
		assert.False(t, HasAtom(targetAtom))
	})
}
