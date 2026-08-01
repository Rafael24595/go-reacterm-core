package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Line func(order bool, lne line.Line) []line.Line
