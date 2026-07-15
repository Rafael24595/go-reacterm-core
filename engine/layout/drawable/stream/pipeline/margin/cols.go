package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func colsDrawTransformer(hint hint.Size[winsize.Cols], opts ...cols.Option) pipeline.DataTransformer {
	transformer := margin.Cols(hint, opts...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []line.Line, hasNext bool) ([]line.Line, bool) {
		return transformer(size, lines), hasNext
	}
}
