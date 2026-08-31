package state

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestUIState_Initialization(t *testing.T) {
	uiState := NewUIState()

	assert.NotNil(t, uiState.Store)
	assert.False(t, uiState.Helper.ShowHelp)
	assert.Equal(t, 0, uiState.Pager.CurrentPage)
}
