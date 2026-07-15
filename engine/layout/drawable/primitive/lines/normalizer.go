package lines

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

type linesNormalizer func() []wrap.LayoutLine

func eagerNormalizer(lines ...wrap.LayoutLine) linesNormalizer {
	return func() []wrap.LayoutLine {
		return lines
	}
}

func lazyNormalizer(lines ...line.Line) linesNormalizer {
	return func() []wrap.LayoutLine {
		return wrap.NormalizeLines(lines...)
	}
}
