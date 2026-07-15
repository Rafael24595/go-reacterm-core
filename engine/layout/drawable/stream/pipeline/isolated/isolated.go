package isolated

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/wipe"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "isolated_pipeline"

func Unit(unit drawable.Unit) drawable.Unit {
	unt := pipeline.New(unit).
		PushBootSteps(wipe.BootTransformer()).
		SetDrawStep(drain.DrawTransformer(true)).
		ToUnit()

	unt.Name = Name
	return unt
}

func UnitFromLines(lns ...text.Line) drawable.Unit {
	return Unit(
		lines.FromLines(lns...).ToUnit(),
	)
}
