package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// Line defines a function strategy for processing and transforming a line into one or more lines based on rendering direction.
type Line func(order bool, lne line.Line) []line.Line
