package styler

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

func TestAtomRegistryCoverage(t *testing.T) {
	for r := range atom.Registry() {
		var found bool

		for _, s := range atoms {
			if s.Atom == r.Atom() {
				found = true
				break
			}
		}

		assert.True(t, found, "atom %s not found", r.Name)
	}
}

func TestAtomsAreUnique(t *testing.T) {
	cache := make(map[atom.Atom]bool)

	for _, v := range atoms {
		_, ok := cache[v.Atom]
		assert.False(t, ok)

		cache[v.Atom] = true
	}
}

func TestToLower(t *testing.T) {
	assert.Equal(t, "hello world", toLower("Hello World"))
	assert.Equal(t, "", toLower(""))
}

func TestToUpper(t *testing.T) {
	assert.Equal(t, "HELLO WORLD", toUpper("Hello World"))
	assert.Equal(t, "", toUpper(""))
}

func TestToDefault(t *testing.T) {
	assert.Equal(t, "Hello World", toDefault("Hello World"))
	assert.Equal(t, "", toDefault(""))
}
