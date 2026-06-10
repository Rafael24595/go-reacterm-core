package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Rows(hintSize hint.Size[winsize.Rows], opts ...rows.Option) transform.Transformer {
	cfg := rows.ResolveConfig(opts...)

	return func(size winsize.Winsize, lines []text.Line) []text.Line {
		linesLen := winsize.Rows(len(lines))

		margin := hintSize.Min(size.Rows) * VerticalFactor(cfg.Position)

		transformer := padding.Rows(
			hint.Fixed(linesLen+margin),
			opts...,
		)

		return transformer(size, lines)
	}
}
