package stack

import (
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type LayerRenderer func(winsize.Winsize, drawable.Unit) ([]text.Line, bool)

func defaultRenderer(size winsize.Winsize, unit drawable.Unit) ([]text.Line, bool) {
	return drain.Unit(size, unit, true)
}
