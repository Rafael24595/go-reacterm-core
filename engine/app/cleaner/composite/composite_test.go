package composite

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"

	cleaner_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/cleaner"
)

func TestComposite_ToStateCleaner(t *testing.T) {
	cleaner := NewCleaner()

	cleaner_test.Helper_ToStateCleaner(t, cleaner)
}

func TestCompositeCleanup(t *testing.T) {

	c := NewCleaner(
		func(r screen.Result, s *state.UIState) *state.UIState {
			s.Pager.ActualPage = 0
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
	uiState.Pager.ActualPage = 10
	uiState.Pager.ForceShow = false
	uiState.Helper.ShowHelp = true

	res := screen.Result{}

	uiState = c.Cleanup(res, uiState)

	assert.Equal(t, 0, uiState.Pager.ActualPage)
	assert.True(t, uiState.Pager.ForceShow)
	assert.False(t, uiState.Helper.ShowHelp)
}
