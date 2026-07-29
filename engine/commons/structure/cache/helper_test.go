package cache

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func Helper_ToCache(t *testing.T, cache Cache[any, any]) {
	t.Helper()

	assert.NotNil(t, cache.Get, "Cache.Get should be set")
	assert.NotNil(t, cache.Put, "Cache.Put should be set")
	assert.NotNil(t, cache.Del, "Cache.Del should be set")
	assert.NotNil(t, cache.Len, "Cache.Len should be set")
	assert.NotNil(t, cache.Cls, "Cache.Cls should be set")
}
