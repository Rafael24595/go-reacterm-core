package lines

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

const separator = " | "

func NextIndexedLine(
	cols winsize.Cols,
	lines []layout.Line,
	meta indexMeta,
) (*line.Line, []layout.Line) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]layout.Line, 0)
	}

	prefix, lines := extractPrefix(lines, meta)

	fixedCols := cols.Sub(meta.totalWidth)

	assert.True(fixedCols > 0, "index prefix should be lesser than line size")

	cursor, rest := wrap.NextBuilder(fixedCols, lines)
	if cursor == nil {
		return nil, rest
	}

	cursor.UnshiftFrags(
		frag.FromString(prefix),
	)

	return cursor.LinePtr(), rest
}

func extractPrefix(
	lines []layout.Line,
	meta indexMeta,
) (string, []layout.Line) {
	if lines[0].Source.Order() == 0 {
		return meta.body(), lines

	}
	order := int(lines[0].Source.Order())
	prefix := meta.header(order)

	lines[0].Source = line.BuilderFromLine(lines[0].Source).
		SetOrder(0).
		Line()

	return prefix, lines
}

func computeIndexMeta(lines []layout.Line) *indexMeta {
	size := winsize.Cols(0)

	for _, line := range lines {
		if line.Source.Order() == 0 {
			continue
		}

		digits := math.Digits(line.Source.Order())
		size = max(size, winsize.Cols(digits))
	}

	if size == 0 {
		return nil
	}

	prefix := format.PatternRight(
		size, format.TextFromString(marker.DefaultPaddingText),
	)

	return &indexMeta{
		sufix:      separator,
		prefixBody: prefix,
		digits:     uint16(size),
		totalWidth: size + runes.Measure(separator),
	}
}
