package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const (
	SystemMetaTag = "system_meta"
)

type BootFunc func(state.UIState)
type KeysFunc func() Definition
type TickFunc func(*state.UIState, Event) Result
type ViewFunc func(state.UIState) viewmodel.ViewModel

type Funcs interface {
	BootFunc | KeysFunc | TickFunc | ViewFunc
}

type Screen struct {
	Boot BootFunc
	Keys KeysFunc
	Tick TickFunc
	View ViewFunc
}

func IsZeroScreen(screen Screen) bool {
	return screen.Boot == nil ||
		screen.Keys == nil ||
		screen.Tick == nil ||
		screen.View == nil
}
