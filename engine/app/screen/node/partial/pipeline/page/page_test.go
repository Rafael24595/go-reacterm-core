package page

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"

	pager_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/pager"
)

func TestUse(t *testing.T) {
	code := action.Kind(10)
	mock := &pager_test.MockStrategy{
		ActionKind: code,
	}

	transformer := Use(mock.ToStrategy().Action)

	vm := transformer(viewmodel.ViewModel{
		Pager: pager.NewStrategy().
			SetAction(action.Paged()),
	})

	assert.Equal(t, code, vm.Pager.Action.Kind)
	assert.NotEqual(t, action.Paged().Kind, vm.Pager.Action.Kind)
}
