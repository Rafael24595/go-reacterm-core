package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/margin"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func rowsDataTransformer(hint hint.Size[winsize.Rows], opts ...rows.Option) pipeline.DataTransformer {
	transformer := margin.Rows(hint, opts...)
	return func(size winsize.Winsize, _ drawable.Unit, lines []line.Line, hasNext bool) ([]line.Line, bool) {
		return transformer(size, lines), hasNext
	}
}
