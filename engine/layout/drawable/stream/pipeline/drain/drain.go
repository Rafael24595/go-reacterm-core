package drain

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
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

func UnitFromLines(lines ...text.Line) drawable.Unit {
	return Unit(
		line.FromLines(lines...).ToUnit(),
	)
}

func UnitFromFrags(frags ...text.Frag) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFrags(frags...),
	)
}

func UnitFromString(txt ...string) drawable.Unit {
	return UnitFromLines(
		*text.LineFromFrags(
			text.FragsFromString(txt...)...,
		),
	)
}
