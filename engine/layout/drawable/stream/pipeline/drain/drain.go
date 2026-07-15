package drain

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

const Name = "drain_pipeline"

func DrawTransformer(lazy bool) pipeline.DrawTransformer {
	return func(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
		return drain.Unit(size, unit, lazy)
	}
}

func Unit(unit drawable.Unit) drawable.Unit {
	unt := pipeline.New(unit).
		SetDrawStep(DrawTransformer(true)).
		ToUnit()

	unt.Name = Name
	return unt
}

func UnitFromLines(lns ...text.Line) drawable.Unit {
	return Unit(
		lines.FromLines(lns...).ToUnit(),
	)
}

func UnitFromFrags(frs ...frag.Frag) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFrags(frs...),
	)
}

func UnitFromString(txt ...string) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFrags(
			frag.FromStrings(txt...)...,
		),
	)
}
