package modalmenu

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestModalMenu_ToNode(t *testing.T) {
	node := New().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, "base")
}

func TestIndexMenu_Boot(t *testing.T) {
	menu := New().
		AddOptions(
			input.MenuOption{Id: "4"},
			input.MenuOption{Id: "3"},
			input.MenuOption{Id: "2"},
			input.MenuOption{Id: "1"},
		)
	node := menu.ToNode()

	uiState := state.NewUIState()
	KeySync.Set(
		uiState.Store,
		node.Name,
		"1",
	)

	node.Screen.Boot(*uiState)

	_, ok := KeySync.Get(uiState.Store, menu.reference)

	assert.False(t, ok)

	assert.Equal(t, 3, menu.cursor)
}

func TestModalMenu_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
