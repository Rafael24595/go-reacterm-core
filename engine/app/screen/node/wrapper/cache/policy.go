package cache

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
)

type Policy func(screen.Result, cache.Cache[hash.Hash, delta.Delta])

func CleanOnNewNode(result screen.Result, cache cache.Cache[hash.Hash, delta.Delta]) {
	if result.HasNode() {
		cache.Cls()
	}
}
