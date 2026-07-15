package transform

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Transformer func(size winsize.Winsize, lines []line.Line) []line.Line
