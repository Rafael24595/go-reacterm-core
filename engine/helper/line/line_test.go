package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

func TestFindLineStartAndEnd(t *testing.T) {
	buf := []rune{'h', 'e', 'l', 'l', 'o', ascii.ENTER_LF, 'g', 'o', 'l', 'a', 'n', 'g'}

	t.Run("FindLineStart at middle of second line", func(t *testing.T) {
		start := FindLineStart(buf, offset.Offset(8))
		assert.Equal(t, offset.Offset(6), start)
	})

	t.Run("FindLineEnd for first line", func(t *testing.T) {
		end := FindLineEnd(buf, offset.Offset(0))
		assert.Equal(t, offset.Offset(5), end)
	})

	t.Run("DistanceFromLF", func(t *testing.T) {
		dist := DistanceFromLF(buf, offset.Offset(8))
		assert.Equal(t, offset.Offset(2), dist)
	})
}

func TestFindNextAndPrevLineStart(t *testing.T) {
	buf := []rune("" +
		"golang" +
		string(rune(ascii.ENTER_LF)) +
		"rust" +
		string(rune(ascii.ENTER_LF)) +
		"zig",
	)

	t.Run("FindNextLineStart", func(t *testing.T) {
		nextStart, ok := FindNextLineStart(buf, offset.Offset(2))
		assert.True(t, ok)
		assert.Equal(t, offset.Offset(7), nextStart)
	})

	t.Run("FindPrevLineStart", func(t *testing.T) {
		prevStart, ok := FindPrevLineStart(buf, offset.Offset(8))
		assert.True(t, ok)
		assert.Equal(t, offset.Offset(0), prevStart)
	})
}

func TestClampToLine(t *testing.T) {
	buf := []rune(
		"go" + string(rune(ascii.ENTER_LF)) + "lang",
	)

	t.Run("col within bounds", func(t *testing.T) {
		clamped := ClampToLine(buf, offset.Offset(0), offset.Offset(2))
		assert.Equal(t, offset.Offset(2), clamped)
	})

	t.Run("col exceeds line length", func(t *testing.T) {
		clamped := ClampToLine(buf, offset.Offset(0), offset.Offset(10))
		assert.Equal(t, offset.Offset(2), clamped)
	})
}
