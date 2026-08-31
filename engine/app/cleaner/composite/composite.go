package composite

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func Cleaner(cls ...cleaner.Cleaner) cleaner.Cleaner {
	return func(res screen.Result, uiState *state.UIState) *state.UIState {
		for _, part := range cls {
			uiState = part(res, uiState)
		}
		return uiState
	}
}
