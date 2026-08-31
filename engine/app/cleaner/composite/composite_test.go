package composite

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestComposite(t *testing.T) {
	cleaner := Cleaner(
		func(r screen.Result, s *state.UIState) *state.UIState {
			s.Pager.CurrentPage = 0
			return s
		},
		func(r screen.Result, s *state.UIState) *state.UIState {
			s.Pager.ForceShow = true
			return s
		},
		func(r screen.Result, s *state.UIState) *state.UIState {
			s.Helper.ShowHelp = false
			return s
		},
	)

	uiState := state.NewUIState()
	uiState.Pager.CurrentPage = 10
	uiState.Pager.ForceShow = false
	uiState.Helper.ShowHelp = true

	res := screen.Result{}

	uiState = cleaner(res, uiState)

	assert.Equal(t, 0, uiState.Pager.CurrentPage)
	assert.True(t, uiState.Pager.ForceShow)
	assert.False(t, uiState.Helper.ShowHelp)
}

func TestComposite_EmptyListIsNoOp(t *testing.T) {
	cleaner := Cleaner()

	uiState := state.NewUIState()
	uiState.Pager.CurrentPage = 5

	result := cleaner(screen.Result{}, uiState)

	assert.Equal(t, uiState, result)
	assert.Equal(t, 5, result.Pager.CurrentPage)
}
