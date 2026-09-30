package splitter

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
)

func TestSplitFragWithCache(t *testing.T) {
	cache := NewFragCache()

	splitterWithCache := SplitFragWithCache(cache)

	frg := frag.FromString("hello world")

	builder1 := delta.NewBuilder().
		WithDelta(splitterWithCache(frg))

	assert.Size(t, 3, builder1.Frags)

	cachedDelta, ok := cache.Get(frg.Hash())
	builder2 := delta.NewBuilder().WithDelta(cachedDelta)

	assert.True(t, ok)
	assert.Equal(t,
		builder1.Frags[0].Base.Text(), builder2.Frags[0].Base.Text(),
	)

	builder3 := delta.NewBuilder().
		WithDelta(splitterWithCache(frg))

	assert.Equal(t, builder1.LeftEdge, builder3.LeftEdge)
	assert.Equal(t, builder1.RightEdge, builder3.RightEdge)
	assert.Size(t, len(builder1.Frags), builder3.Frags)
}

func TestSplitLineWithCache(t *testing.T) {
	c := NewFragCache()
	lineSplitter := SplitLineWithCache(c)

	lne := line.FromString("golang rust zig")

	words, frags := lineSplitter(lne)

	assert.Size(t, 5, frags)
	assert.Size(t, 5, words)

	for f := range lne.All() {
		_, ok := c.Get(f.Hash())
		assert.True(t, ok)
	}
}

func BenchmarkSplitWordsCached(b *testing.B) {
	splitter := SplitLineWithCache(
		NewFragCache(),
	)

	l := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		splitter(l)
	}
}

func BenchmarkSplitWords_ColdCache(b *testing.B) {
	cache := NewFragCache()

	splitter := SplitLineWithCache(cache)
	lne := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		cache.Cls()
		_, _ = splitter(lne)
	}
}

func BenchmarkSplitWords_WarmCache(b *testing.B) {
	cache := NewFragCache()

	splitter := SplitLineWithCache(cache)
	lne := benchmarkLine(2000)

	splitter(lne)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = splitter(lne)
	}
}

func BenchmarkSplitLineWith_WarmCache(b *testing.B) {
	cache := NewFragCache()
	splitter := SplitFragWithCache(cache)

	lne := benchmarkLine(2000)

	for _, frg := range lne.Slice() {
		splitter(frg)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, _ = SplitLineWith(splitter, lne)
	}
}
