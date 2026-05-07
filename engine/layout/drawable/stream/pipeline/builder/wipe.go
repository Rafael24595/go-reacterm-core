package builder

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/transformer/wipe"
)

const NameWipe = "wipe_pipeline"

func WipeFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	drw := pipeline.New(drawable).
		SetDrawStep(wipe.DrawTransformer()).
		ToDrawable()

	drw.Name = NameWipe
	return drw
}
