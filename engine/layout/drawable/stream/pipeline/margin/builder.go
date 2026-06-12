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

func (b *Builder) Rows(hint *hint.Size[winsize.Rows], opts ...rows.Option) *Builder {
	b.hintRows = hint
	b.optionsRows = append(b.optionsRows, opts...)
	return b
}

func (b *Builder) Cols(hint *hint.Size[winsize.Cols], opts ...cols.Option) *Builder {
	b.hintCols = hint
	b.optionsCols = append(b.optionsCols, opts...)
	return b
}

func (b *Builder) Steps() (pipeline.DrawTransformer, []pipeline.DataTransformer) {
	draw := func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		cfgY := rows.ResolveConfig(b.optionsRows...)
		cfgX := cols.ResolveConfig(b.optionsCols...)

		marginY := winsize.Rows(0)
		if b.hintRows != nil {
			marginY = b.hintRows.Min(size.Rows) * margin.VerticalFactor(cfgY.Position)
		}

		marginX := winsize.Cols(0)
		if b.hintCols != nil {
			marginX = b.hintCols.Min(size.Cols) * margin.HorizontalFactor(cfgX.Position)
		}

		fixedSize := winsize.New(
			size.Rows.Sub(marginY),
			size.Cols.Sub(marginX),
		)

		return unit.Drawable.Draw(fixedSize)
	}

	data := make([]pipeline.DataTransformer, 0, 2)

	if b.hintRows != nil {
		data = append(data,
			rowsDataTransformer(*b.hintRows, b.optionsRows...),
		)
	}

	if b.hintCols != nil {
		data = append(data,
			colsDrawTransformer(*b.hintCols, b.optionsCols...),
		)
	}

	return draw, data
}

func (b *Builder) ToUnit(unit drawable.Unit) drawable.Unit {
	draw, data := b.Steps()
	return pipeline.New(unit).
		SetDrawStep(draw).
		PushDataSteps(data...).
		ToUnit()
}
