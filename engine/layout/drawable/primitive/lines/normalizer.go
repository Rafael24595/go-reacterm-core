package lines

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

type linesNormalizer func() []layout.Line

func eagerNormalizer(lines ...layout.Line) linesNormalizer {
	return func() []layout.Line {
		return lines
	}
}

func lazyNormalizer(lines ...line.Line) linesNormalizer {
	return func() []layout.Line {
		return wrap.NormalizeLines(lines...)
	}
}
