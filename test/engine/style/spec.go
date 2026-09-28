package style_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func SpecEquals(t *testing.T, a, b spec.Spec) {
	t.Helper()

	assert.Equal(t, a.Kind(), b.Kind())
	assert.Equal(t, a.Hash(), b.Hash())
	assert.Size(t, len(a.Args()), b.Args())
}
