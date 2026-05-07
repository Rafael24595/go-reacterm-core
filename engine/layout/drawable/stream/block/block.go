package block

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/line"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "block_drawable"

//TODO: Reduce bureaucracy?
type BlockDrawable struct {
	lazy     bool
	drawable drawable.Drawable
}

func New(drawable drawable.Drawable) *BlockDrawable {
	return &BlockDrawable{
		lazy:     true,
		drawable: drawable,
	}
}

func DrawableFromDrawable(drawable drawable.Drawable) drawable.Drawable {
	return New(drawable).ToDrawable()
}

func DrawableFromLines(lines ...text.Line) drawable.Drawable {
	return DrawableFromDrawable(
		line.New(lines...).ToDrawable(),
	)
}

func DrawableFromString(txt ...string) drawable.Drawable {
	lines := text.LineFromFragments(
		text.FragmentsFromString(txt...)...,
	)
	return DrawableFromLines(*lines)
}

func (d *BlockDrawable) Lazy(lazy bool) *BlockDrawable {
	d.lazy = lazy
	return d
}

func (d *BlockDrawable) ToDrawable() drawable.Drawable {
	drw := pipeline.New(d.drawable).
		SetDrawStep(drain.DrawTransformer(d.lazy)).
		ToDrawable()

	drw.Name = Name
	return drw
}
