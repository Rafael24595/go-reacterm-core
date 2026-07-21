package lines

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestWrapNextLine_FitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: "    ",
	}

	layout := wrap.NormalizeLines(
		line.TextOrdered(1, "golang"),
	)

	got, remain := NextIndexedLine(10, layout, meta)

	assert.Equal(t, "1 | golang", text_test.LineToString(*got))

	assert.Empty(t, remain)
}

func TestWrapNextLine_SplitWithMeta(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	layout := wrap.NormalizeLines(
		line.TextOrdered(1, "golang rust"),
	)

	got, remain := NextIndexedLine(10, layout, meta)

	assert.Equal(t, "1 | golang", text_test.LineToString(*got))
	assert.Size(t, 1, remain)

	got, remain = NextIndexedLine(10, remain, meta)

	assert.Equal(t, "  |  rust", text_test.LineToString(*got))
	assert.Empty(t, remain)
}

func TestWrapNextLine_IndexShouldBeLesser(t *testing.T) {
	meta := indexMeta{
		totalWidth: 4,
		sufix:      " | ",
		digits:     1,
		prefixBody: " ",
	}

	layout := wrap.LayoutLine{
		Source: line.TextOrdered(1, "golang"),
	}

	assert.Panic(t, func() {
		NextIndexedLine(4, []wrap.LayoutLine{layout}, meta)
	})
}
