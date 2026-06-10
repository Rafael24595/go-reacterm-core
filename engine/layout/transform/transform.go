package transform

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Transformer func(size winsize.Winsize, lines []text.Line) []text.Line
