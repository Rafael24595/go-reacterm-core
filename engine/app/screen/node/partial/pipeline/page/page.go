package page

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

func Use(engine pager.Engine) pipeline.Transformer {
	return func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		vm.Pager.Engine = engine
		return vm
	}
}
