package gutter

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

const Name = "gutter_pipeline"

func DrawTransformer(opts ...Option) pipeline.DrawTransformer {
	meta := newMeta(opts...)

	leftMeasure := runes.Measure(meta.left)
	rightMeasure := runes.Measure(meta.right)
	measure := leftMeasure + rightMeasure

	leftFrg := frag.New(meta.left)
	rightFrg := frag.New(meta.right)

	return func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		if measure >= size.Cols {
			return unit.Drawable.Draw(size)
		}

		fixedSize := winsize.New(
			size.Rows,
			size.Cols.Sub(measure),
		)

		lines, hasNext := unit.Drawable.Draw(fixedSize)
		for i := range lines {
			if leftMeasure > 0 {
				lines[i].UnshiftFrags(*leftFrg)
			}
			if rightMeasure > 0 {
				lines[i].PushFrags(*rightFrg)
			}
		}

		return lines, hasNext
	}
}

func Unit(unit drawable.Unit, opts ...Option) drawable.Unit {
	unt := pipeline.New(unit).
		SetDrawStep(
			DrawTransformer(opts...),
		).
		ToUnit()

	unt.Name = Name
	return unt
}
