package screen

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const (
	// SystemMetaTag defines the metadata tag used for system-level operations.
	SystemMetaTag = "system_meta"
)

// BootFunc executes the screen preloading logic prior to rendering.
type BootFunc func(state.UIState)

// KeysFunc returns the required key bindings and shortcut definitions for the screen.
type KeysFunc func() Definition

// TickFunc processes user input events and executes the associated business logic for the screen.
type TickFunc func(*state.UIState, Event) Result

// ViewFunc returns the required resources (ViewModel) to render the screen interface.
type ViewFunc func(state.UIState) viewmodel.ViewModel

// Funcs defines a constraint for all valid screen function types.
type Funcs interface {
	BootFunc | KeysFunc | TickFunc | ViewFunc
}

// Screen represents a UI screen component along with its lifecycle functions.
type Screen struct {
	// Boot executes the screen preloading logic.
	Boot BootFunc
	// Keys returns the key definitions required by the screen.
	Keys KeysFunc
	// Tick handles input events and updates the screen's state and business logic.
	Tick TickFunc
	// View provides the ViewModel required for screen rendering.
	View ViewFunc
}

// IsZeroScreen checks whether any required lifecycle function in the screen is uninitialized (nil).
func IsZeroScreen(screen Screen) bool {
	return screen.Boot == nil ||
		screen.Keys == nil ||
		screen.Tick == nil ||
		screen.View == nil
}
