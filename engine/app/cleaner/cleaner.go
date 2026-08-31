package cleaner

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

// Cleaner defines a function contract tasked with processing a screen result and returning
// a reconciled or pruned UIState.
type Cleaner func(screen.Result, *state.UIState) *state.UIState
