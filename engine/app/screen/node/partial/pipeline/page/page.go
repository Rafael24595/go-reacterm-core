package page

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func Use(action action.Action) pipeline.Transformer {
	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		vm.Pager.Action = action
		return vm
	}
}
