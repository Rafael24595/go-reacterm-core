package draw

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type DrawContext struct {
	State *state.UIState
	Size  winsize.Winsize
}

func NewDrawContext(uiState *state.UIState, size winsize.Winsize) *DrawContext {
	return &DrawContext{
		State: uiState,
		Size:  size,
	}
}
