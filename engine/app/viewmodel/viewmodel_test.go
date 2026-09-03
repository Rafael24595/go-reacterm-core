package viewmodel_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func TestViewModel_New(t *testing.T) {
	vm := viewmodel.New()

	assert.NotNil(t, vm.Header)
	assert.NotNil(t, vm.Kernel)
	assert.NotNil(t, vm.Footer)
	assert.NotNil(t, vm.Pager)
}

func TestViewModel_CloneIsIsolated(t *testing.T) {
	original := viewmodel.New()

	cloned := original.Clone()

	assert.NotNil(t, cloned)
	assert.NotSame(t, original.Header, cloned.Header)
	assert.NotSame(t, original.Kernel, cloned.Kernel)
	assert.NotSame(t, original.Footer, cloned.Footer)
}
