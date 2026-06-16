package help

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestHelp_ToNode(t *testing.T) {
	name := "base"
	mock := screen_test.MockNode{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, name)
}

func TestHelp_Propagate(t *testing.T) {
	name := "base"
	mock := screen_test.MockNode{
		Name: name,
	}

	node := New(mock.ToNode()).ToNode()
	screen_test.Helper_Propagate(t, name, 0, node)
}

func TestHelp_ToggleHelpKey(t *testing.T) {
	called := false

	mock := screen_test.MockNode{}

	node := New(mock.ToNode()).ToNode()

	state := &state.UIState{}
	event := screen.Event{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	}

	node.Screen.Tick(state, event)

	assert.True(t, state.Helper.ShowHelp)
	assert.False(t, called)
}

func TestHelp_DelegatesTickWhenKeyRequired(t *testing.T) {
	called := false

	action := key.CustomActionHelp
	definition := screen.DefinitionFromActions(action)

	mock := screen_test.MockNode{
		Definition: &definition,
		Tick: func(s *state.UIState, e screen.Event) screen.Result {
			called = true
			return screen.EmptyResult()
		},
	}

	node := New(mock.ToNode()).ToNode()

	state := &state.UIState{}
	event := screen.Event{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	}

	node.Screen.Tick(state, event)

	assert.False(t, state.Helper.ShowHelp)
	assert.True(t, called)
}

func TestHelp_WrapsReturnedScreen(t *testing.T) {
	called := false

	action := key.ActionEnter
	definition := screen.DefinitionFromActions(action)

	mockNext := screen_test.MockNode{
		Name: "next",
	}

	mockBase := screen_test.MockNode{
		Definition: &definition,
		Tick: func(s *state.UIState, _ screen.Event) screen.Result {
			called = true
			next := mockNext.ToNode()
			return screen.Result{
				Node: &next,
			}
		},
	}

	help := New(mockBase.ToNode())
	wrapped := help.ToNode()

	uiState := &state.UIState{}
	event := screen.Event{
		Key: *key.NewKeyCode(key.ActionEnter),
	}

	wrapped.Screen.Tick(uiState, screen.Event{
		Key: *key.NewKeyCode(key.CustomActionHelp),
	})

	assert.True(t, uiState.Helper.ShowHelp)

	result := wrapped.Screen.Tick(uiState, event)

	assert.True(t, called)
	assert.NotNil(t, result.Node)
	assert.Equal(t, "next", result.Node.Name)
}
