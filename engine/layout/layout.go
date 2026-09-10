package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Composer func(*state.UIState, winsize.Winsize, viewmodel.ViewModel) (*state.UIState, []line.Line)

type Layout struct {
	Compose Composer
}
