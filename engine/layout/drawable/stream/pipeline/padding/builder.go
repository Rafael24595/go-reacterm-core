package padding

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Builder struct {
	hintRows    *hint.Size[winsize.Rows]
	optionsRows []rows.Option
	hintCols    *hint.Size[winsize.Cols]
	optionsCols []cols.Option
}

func NewBuilder() *Builder {
	return &Builder{
		hintRows:    nil,
		optionsRows: make([]rows.Option, 0),
		hintCols:    nil,
		optionsCols: make([]cols.Option, 0),
	}
}

func (b *Builder) Rows(hint hint.Size[winsize.Rows], opts ...rows.Option) *Builder {
	b.hintRows = &hint
	b.optionsRows = append(b.optionsRows, opts...)
	return b
}

func (b *Builder) Cols(hint hint.Size[winsize.Cols], opts ...cols.Option) *Builder {
	b.hintCols = &hint
	b.optionsCols = append(b.optionsCols, opts...)
	return b
}

func (b *Builder) Steps() []pipeline.DataTransformer {
	data := make([]pipeline.DataTransformer, 0, 2)

	if b.hintRows != nil {
		data = append(data,
			Rows(*b.hintRows, b.optionsRows...),
		)
	}

	if b.hintCols != nil {
		data = append(data,
			Cols(*b.hintCols, b.optionsCols...),
		)
	}

	return data
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	return pipeline.New(unit).
		PushDataSteps(b.Steps()...).
		ToUnit()
}
