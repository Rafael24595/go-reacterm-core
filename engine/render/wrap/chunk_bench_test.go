package wrap

import (
	"fmt"
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
)

func BenchmarkChunkSizes(b *testing.B) {
	sizes := []offset.Offset{
		16,
		32,
		64,
		128,
		256,
	}

	src := benchmarkLine(100_000)

	for _, size := range sizes {
		b.Run(fmt.Sprintf("chunk_%d", size), func(b *testing.B) {
			for b.Loop() {
				_ = chunk.Line(src, size)
			}
		})
	}
}
