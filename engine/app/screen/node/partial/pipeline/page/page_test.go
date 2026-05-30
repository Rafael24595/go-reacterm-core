package page

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	pager_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/pager"
)

func TestUse(t *testing.T) {
	code := pager.EngineCode(10)
	mock := &pager_test.MockStrategy{
		EngineCode: code,
	}

	transformer := Use(mock.ToStrategy().Engine)

	vm := transformer(viewmodel.ViewModel{
		Pager: pager.NewStrategy().
			SetEngine(pager.EnginePage()),
	})

	assert.Equal(t, code, vm.Pager.Engine.Code)
	assert.NotEqual(t, pager.EnginePage().Code, vm.Pager.Engine.Code)
}
