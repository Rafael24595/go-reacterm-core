package wipe

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const Name = "wipe_pipeline"

func BootTransformer() pipeline.BootTransformer {
	return func(size winsize.Winsize, unit drawable.Unit) drawable.Unit {
		unit.Drawable.Wipe()
		return unit
	}
}

func DrawTransformer() pipeline.DrawTransformer {
	transformer := DataTransformer()
	return func(size winsize.Winsize, unit drawable.Unit) ([]line.Line, bool) {
		lines, hasNext := unit.Drawable.Draw(size)
		return transformer(size, unit, lines, hasNext)
	}
}

func DataTransformer() pipeline.DataTransformer {
	return func(_ winsize.Winsize, unit drawable.Unit, lines []line.Line, hasNext bool) ([]line.Line, bool) {
		if len(lines) == 0 {
			return lines, hasNext
		}

		if !hasNext {
			unit.Drawable.Wipe()
		}

		return lines, true
	}
}

func Drawable(unit drawable.Unit) drawable.Unit {
	unt := pipeline.New(unit).
		SetDrawStep(DrawTransformer()).
		ToUnit()

	unt.Name = Name
	return unt
}
