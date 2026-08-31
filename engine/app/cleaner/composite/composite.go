package composite

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

// Cleaner combines multiple cleaner functions into a single Cleaner execution chain,
// executing them sequentially in the order they were provided.
func Cleaner(cls ...cleaner.Cleaner) cleaner.Cleaner {
	return func(res screen.Result, uiState *state.UIState) *state.UIState {
		for _, part := range cls {
			uiState = part(res, uiState)
		}
		return uiState
	}
}
