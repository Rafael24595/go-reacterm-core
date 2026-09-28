package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestLookupDescriptor(t *testing.T) {
	t.Run("returns descriptor with valid name for JustifyRight", func(t *testing.T) {
		desc, ok := Lookup(KindJustifyRight)

		assert.True(t, ok)
		assert.Equal(t, KindJustifyRight, desc.Kind())
		assert.Equal(t, "JustifyRight", desc.Name())

		var keys []ArgKey
		for key := range desc.Args() {
			keys = append(keys, key)
		}
		
		assert.Equal(t, 2, len(keys))
	})

	t.Run("returns false for unknown Kind", func(t *testing.T) {
		unregistered := Kind(1 << 30)
		_, ok := Lookup(unregistered)
		assert.False(t, ok)
	})
}

func TestDescriptorRegistryIterator(t *testing.T) {
	var kinds []Kind
	for desc := range Registry() {
		kinds = append(kinds, desc.Kind())
	}

	assert.True(t, len(kinds) >= 8)
	assert.Equal(t, KindJustifyRight, kinds[0])
	assert.Equal(t, KindJustifyLeft, kinds[1])
}