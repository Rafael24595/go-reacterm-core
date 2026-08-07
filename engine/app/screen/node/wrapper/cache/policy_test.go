package cache

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"

	structure_test "github.com/Rafael24595/go-reacterm-core/test/engine/commons/structure"
)

func TestCleanOnNewNode_Concrete(t *testing.T) {
	t.Run("Clean the cache when node is not nil", func(t *testing.T) {
		mock := structure_test.MockFromCache(
			cache.NewMemory[hash.Hash, delta.Delta](),
		)

		result := screen.ResultFromNode(
			screen.Node{},
		)

		CleanOnNewNode(result, mock.ToCache())

		assert.GreaterThan(t, 0, mock.ClsCalls)
	})

	t.Run("Keep the cache when node is nil", func(t *testing.T) {
		mock := structure_test.MockFromCache(
			cache.NewMemory[hash.Hash, delta.Delta](),
		)

		result := screen.Result{}

		CleanOnNewNode(result, mock.ToCache())

		assert.Equal(t, 0, mock.ClsCalls)
	})
}
