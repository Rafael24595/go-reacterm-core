package processor

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func Chunk(limit offset.Offset) Line {
	return func(_ bool, lne line.Line) []line.Line {
		return []line.Line{
			chunk.Line(lne, limit),
		}
	}
}
