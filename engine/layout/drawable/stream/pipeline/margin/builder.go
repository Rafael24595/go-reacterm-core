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
	hintY    *hint.Size[winsize.Rows]
	optionsY []rows.Option
	hintX    *hint.Size[winsize.Cols]
	optionsX []cols.Option
}

func NewBuilder() *Builder {
	return &Builder{
		hintY:    nil,
		optionsY: make([]rows.Option, 0),
		hintX:    nil,
		optionsX: make([]cols.Option, 0),
	}
}

func (b *Builder) Y(hintY *hint.Size[winsize.Rows], opts ...rows.Option) *Builder {
	b.hintY = hintY
	b.optionsY = append(b.optionsY, opts...)
	return b
}

func (b *Builder) X(hintX *hint.Size[winsize.Cols], opts ...cols.Option) *Builder {
	b.hintX = hintX
	b.optionsX = append(b.optionsX, opts...)
	return b
}

func (b *Builder) Steps() (pipeline.DrawTransformer, []pipeline.DataTransformer) {
	draw := func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		cfgY := rows.ResolveConfig(b.optionsY...)
		cfgX := cols.ResolveConfig(b.optionsX...)

		marginY := winsize.Rows(0)
		if b.hintY != nil {
			marginY = b.hintY.Min(size.Rows) * margin.VerticalFactor(cfgY.Position)
		}

		marginX := winsize.Cols(0)
		if b.hintX != nil {
			marginX = b.hintX.Min(size.Cols) * margin.HorizontalFactor(cfgX.Position)
		}

		fixedSize := winsize.New(
			size.Rows.Sub(marginY),
			size.Cols.Sub(marginX),
		)

		return unit.Drawable.Draw(fixedSize)
	}

	data := make([]pipeline.DataTransformer, 0, 2)

	if b.hintY != nil {
		data = append(data,
			rowsDataTransformer(*b.hintY, b.optionsY...),
		)
	}

	if b.hintX != nil{
		data = append(data,
			colsDrawTransformer(*b.hintX, b.optionsX...),
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
