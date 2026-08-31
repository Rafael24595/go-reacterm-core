package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

// Cleaner inspects the screen transition result and retains only the store entries
// that belong to the active node hierarchy stack.
func Cleaner(result screen.Result, uiState *state.UIState) *state.UIState {
	if node, hasNode := result.TryGetNode(); hasNode {
		uiState.Store.RetainOnly(node.Stack)
	}

	return uiState
}
