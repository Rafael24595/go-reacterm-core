package splitter

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

type Line func(lne line.Line) ([]layout.Word, []layout.Frag)
type Frag func(frg frag.Frag) delta.Delta
