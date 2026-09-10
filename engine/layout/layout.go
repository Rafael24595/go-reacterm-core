package layout

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// Composer is a function type that takes the current UIState, window dimensions, and ViewModel,
// returning the updated UIState and the resulting slice of rendered lines.
type Composer func(*state.UIState, winsize.Winsize, viewmodel.ViewModel) (*state.UIState, []line.Line)

// Layout encapsulates the UI composition logic.
type Layout struct {
	// Compose is the function responsible for composing the UI based on the current state, window size, and view model.
	Compose Composer
}
