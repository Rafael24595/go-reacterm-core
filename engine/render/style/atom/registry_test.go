package atom_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

func TestLookup(t *testing.T) {
	t.Run("returns descriptor for registered atom", func(t *testing.T) {
		desc, ok := atom.Lookup(atom.Bold)
		assert.True(t, ok)
		assert.Equal(t, atom.Bold, desc.Atom())
		assert.Equal(t, "Bold", desc.Name())
	})

	t.Run("returns false for unregistered atom", func(t *testing.T) {
		unregistered := atom.Atom(1 << 15)
		_, ok := atom.Lookup(unregistered)
		assert.False(t, ok)
	})
}

func TestRegistry(t *testing.T) {
	var names []string
	for desc := range atom.Registry() {
		names = append(names, desc.Name())
	}

	assert.GreaterThan(t, 0, names)
	
	assert.Equal(t, "Bold", names[0])
	assert.Equal(t, "Dim", names[1])
}
