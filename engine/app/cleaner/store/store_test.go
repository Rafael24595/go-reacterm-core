package store

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestStore_PreservesStateWhenNoNodeInResult(t *testing.T) {
	uiState := state.NewUIState()
	nodeBase := screen_test.MockByName("base")

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	result := screen.ResultFromUIState(uiState)

	Cleaner(result, uiState)

	value, exists := uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.True(t, exists)
	assert.Equal(t, "golang", value.Text())
}

func TestStore_PreservesActiveState(t *testing.T) {
	uiState := state.NewUIState()

	nodeBase := screen_test.MockByName("base")

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockNode{
		Stack: nodeBase.Stack,
	}.ToNode()

	result := screen.ResultFromUIState(uiState)
	result.SetNode(nodeWrapper)

	Cleaner(result, uiState)

	value, exists := uiState.Store.Find(nodeBase.Name, "lang-1")

	assert.True(t, exists)
	assert.Equal(t, "golang", value.Text())
}

func TestStore_RemovesInactiveState(t *testing.T) {
	uiState := state.NewUIState()

	nodeBase := screen_test.MockByName("base")
	nodeNext := screen_test.MockByName("next")

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockNode{}.ToNode()
	nodeWrapper.Stack = nodeNext.Stack

	result := screen.ResultFromUIState(uiState)
	result.SetNode(nodeWrapper)

	Cleaner(result, uiState)

	_, exists := uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.False(t, exists)

	uiState.Store.Push(nodeNext.Name, "lang-2", "ziglang")

	value, exists := uiState.Store.Find(nodeNext.Name, "lang-2")
	assert.True(t, exists)
	assert.Equal(t, "ziglang", value.Text())
}

func TestStore_TransitionBetweenScreens(t *testing.T) {
	uiState := state.NewUIState()

	nodeBase := screen_test.MockByName("base")
	nodeNext := screen_test.MockByName("next")

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockNode{}.ToNode()
	nodeWrapper.Stack = nodeBase.Stack

	result := screen.ResultFromUIState(uiState)
	result.SetNode(nodeWrapper)
	Cleaner(result, uiState)

	_, exists := uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.True(t, exists)

	nodeWrapper.Stack = nodeNext.Stack

	result = screen.ResultFromUIState(uiState)
	result.SetNode(nodeWrapper)
	Cleaner(result, uiState)

	_, exists = uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.False(t, exists)
}
