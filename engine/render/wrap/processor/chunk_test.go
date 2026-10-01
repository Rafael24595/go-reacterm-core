package processor

import (
	"fmt"
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/render/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func benchmarkText(size int) string {
	const sample = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. "

	var b strings.Builder
	b.Grow(size)

	for b.Len() < size {
		b.WriteString(sample)
	}

	return b.String()[:size]
}

func benchmarkLine(size int) line.Line {
	return line.FromFrags(
		frag.FromStrings(
			benchmarkText(size),
		)...,
	)
}

func TestChunkProcessor(t *testing.T) {
	inputLine := line.FromString("Hello World Go")
	limit := offset.Offset(5)

	proc := Chunk(limit)
	lines := proc(false, inputLine)

	assert.Equal(t, 3, lines[0].Size())
}

func BenchmarkChunk(b *testing.B) {
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
