package splitter

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

// Line defines a function strategy for breaking a line into layout words and fragments.
type Line func(lne line.Line) ([]layout.Word, []layout.Frag)
// Frag defines a function strategy for calculating delta metrics from a single fragment.
type Frag func(frg frag.Frag) delta.Delta
