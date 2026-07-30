package wrap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NewFragCache() cache.Cache[hash.Hash, Delta] {
	return cache.NewMemory[hash.Hash, Delta]()
}

func CacheLineWords(cache cache.Cache[hash.Hash, Delta]) LineSplitter {
	processor := CacheFragWords(cache)
	return func(line line.Line) ([]word, []wordFrag) {
		return SplitLineWith(processor, line)
	}
}

func CacheFragWords(cache cache.Cache[hash.Hash, Delta]) FragSplitter {
	return func(frg frag.Frag) Delta {
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
