package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestWrapNextLine_FitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: "    ",
	}

	line := text.NewLine("golang").SetOrder(1)

	got, remain := NextIndexedWrappedLine(10, []text.Line{*line}, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))

	assert.Len(t, 0, remain)
}

func TestWrapNextLine_SplitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	line := text.NewLine("golang rust").SetOrder(1)

	got, remain := NextIndexedWrappedLine(10, []text.Line{*line}, meta)

	assert.Equal(t, "1 | golang", text.LineToString(got))
	assert.Len(t, 1, remain)

	got, remain = NextIndexedWrappedLine(10, remain, meta)

	assert.Equal(t, "  |  rust", text.LineToString(got))
	assert.Len(t, 0, remain)
}

func TestWrapNextLine_IndexShouldBeLesser(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	line := text.NewLine("golang").SetOrder(1)

	assert.Panic(t, func() {
		NextIndexedWrappedLine(4, []text.Line{*line}, meta)
	})
}
