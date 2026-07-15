package processor

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func WithPadding(
	transform func(winsize.Winsize) winsize.Winsize,
	inner render.RawProcessor,
) render.Processor {
	filler := marker.DefaultPaddingText

	return func(lines []line.Line, size winsize.Winsize) string {
		r := transform(size)

		rows := min(r.Rows, size.Rows)
		cols := min(r.Cols, size.Cols)
		newSize := winsize.New(rows, cols)

		content := inner(lines, newSize)
		content = normalize(content, rows)

		diffRows := size.Rows.Sub(rows)
		diffCols := size.Cols.Sub(cols)

		topPadding := diffRows / 2
		leftPadding := diffCols / 2

		buffer := make([]string, size.Rows)
		for i := range size.Rows {
			buffer[i] = format.PatternRight(
				size.Cols, format.TextFromString(filler),
			)
		}

		index := topPadding
		for _, line := range content {
			fixed := format.PatternRight(
				leftPadding, format.TextFromString(filler),
			)

			buffer[index] = format.JustifyLeft(
				size.Cols,
				format.TextFromString(fixed+line),
				filler,
			)
			index += 1
		}

		return strings.Join(buffer, "\n")
	}
}

func normalize(lines []string, rows winsize.Rows) []string {
	buffer := make([]string, rows)
	copy(buffer, lines)
	return buffer
}
