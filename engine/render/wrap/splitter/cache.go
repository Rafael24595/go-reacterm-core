package splitter

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

// NewFragCache constructs an in-memory cache instance mapping fragment hashes to calculated layout deltas.
func NewFragCache() cache.Cache[hash.Hash, delta.Delta] {
	return cache.NewMemory[hash.Hash, delta.Delta]()
}

// SplitLineWithCache returns a Line splitting closure that utilizes cached fragment word calculations.
func SplitLineWithCache(cache cache.Cache[hash.Hash, delta.Delta]) Line {
	processor := SplitFragWithCache(cache)
	return func(line line.Line) ([]layout.Word, []layout.Frag) {
		return SplitLineWith(processor, line)
	}
}

// SplitFragWithCache wraps fragment word tokenization with a lookup cache indexed by fragment hash values.
func SplitFragWithCache(cache cache.Cache[hash.Hash, delta.Delta]) Frag {
	return func(frg frag.Frag) delta.Delta {
		if item, ok := cache.Get(frg.Hash()); ok {
			return item
		}

		frgs := SplitFragByWords(frg)

		cache.Put(
			frg.Hash(), frgs,
		)

		return frgs
	}
}
