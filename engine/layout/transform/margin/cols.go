package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func Cols(hintSize hint.Size[winsize.Cols], opts ...cols.Option) transform.Transformer {
	cfg := cols.ResolveConfig(opts...)

	return func(size winsize.Winsize, lines []line.Line) []line.Line {
		newLines := make([]line.Line, len(lines))

		margin := hintSize.Min(size.Cols) * HorizontalFactor(cfg.Position)

		for i := range lines {
			measure := frag.Measure(size.Cols, lines[i].Text...) + margin

			cols := size.Cols + margin
			if cols.Sub(measure) == 0 {
				newLines[i] = lines[i]
				continue
			}

			remaining := margin
			if measure > size.Cols {
				remaining = measure.Sub(size.Cols)
			}

			newLines[i] = padding.AddColsPadding(
				remaining, lines[i], opts...,
			)
		}

		return newLines
	}
}
