package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Rows(rows hint.Size[winsize.Rows], opts ...rows.Option) pipeline.DataTransformer {
	transformer := margin.Rows(rows, opts...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}

func Cols(cols hint.Size[winsize.Cols], opts ...cols.Option) pipeline.DataTransformer {
	transformer := margin.Cols(cols, opts...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []text.Line, hasNext bool) ([]text.Line, bool) {
		return transformer(size, lines), hasNext
	}
}
