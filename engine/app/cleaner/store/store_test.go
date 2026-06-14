package store

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"

	cleaner_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/cleaner"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestStore_ToStateCleaner(t *testing.T) {
	cleaner := NewCleaner()

	cleaner_test.Helper_ToStateCleaner(t, cleaner)
}

func TestStore__PreservesActiveState(t *testing.T) {
	cleaner := NewCleaner()
	uiState := state.NewUIState()

	nodeBase := screen_test.MockNode{
		Name: "base",
	}.ToNode()

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockNode{
		Stack: nodeBase.Stack,
	}.ToNode()

	result := screen.ResultFromUIState(uiState)
	result.Node = &nodeWrapper

	cleaner.Cleanup(result, uiState)

	value, exists := uiState.Store.Find(nodeBase.Name, "lang-1")

	assert.True(t, exists)
	assert.Equal(t, "golang", value.Stringf())
}

func TestStore__RemovesInactiveState(t *testing.T) {
	cleaner := NewCleaner()
	uiState := state.NewUIState()

	nodeBase := screen_test.MockNode{
		Name: "base",
	}.ToNode()

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeNext := screen_test.MockNode{
		Name: "next",
	}.ToNode()

	nodeWrapper := screen_test.MockNode{}.ToNode()
	nodeWrapper.Stack = nodeNext.Stack

	result := screen.ResultFromUIState(uiState)
	result.Node = &nodeWrapper

	cleaner.Cleanup(result, uiState)

	_, exists := uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.False(t, exists)

	uiState.Store.Push(nodeNext.Name, "lang-2", "ziglang")

	value, exists := uiState.Store.Find(nodeNext.Name, "lang-2")
	assert.True(t, exists)
	assert.Equal(t, "ziglang", value.Stringf())
}

func TestStore__TransitionBetweenScreens(t *testing.T) {
	cleaner := NewCleaner()
	uiState := state.NewUIState()

	nodeBase := screen_test.MockNode{
		Name: "base",
	}.ToNode()

	nodeNext := screen_test.MockNode{
		Name: "next",
	}.ToNode()

	uiState.Store.Push(nodeBase.Name, "lang-1", "golang")

	nodeWrapper := screen_test.MockNode{}.ToNode()
	nodeWrapper.Stack = nodeBase.Stack

	result := screen.ResultFromUIState(uiState)
	result.Node = &nodeWrapper
	cleaner.Cleanup(result, uiState)

	_, exists := uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.True(t, exists)

	nodeWrapper.Stack = nodeNext.Stack

	result = screen.ResultFromUIState(uiState)
	result.Node = &nodeWrapper
	cleaner.Cleanup(result, uiState)

	_, exists = uiState.Store.Find(nodeBase.Name, "lang-1")
	assert.False(t, exists)
}
