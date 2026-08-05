package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func NewCleaner() cleaner.StateCleaner {
	return cleaner.StateCleaner{
		Cleanup: Cleanup,
	}
}

func Cleanup(result screen.Result, uiState *state.UIState) *state.UIState {
	if node, hasNode := result.TryGetNode(); hasNode {
		uiState.Store.RetainOnly(node.Stack)
	}

	return uiState
}
