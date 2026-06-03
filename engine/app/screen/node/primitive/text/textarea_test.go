package text

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTextArea_ToNode(t *testing.T) {
	node := NewArea().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)
	assert.Equal(t, node.Name, "base")
}

func TestTextArea_Init(t *testing.T) {
	area := NewArea()
	node := area.ToNode()

	uiState := state.NewUIState()

	caret := offset.Offset(2)
	anchor := offset.Offset(4)

	store.Push(
		uiState.Store,
		area.reference,
		ArgTextInputState,
		State{
			Buffer: []rune("golang"),
			Caret:  &caret,
			Anchor: &anchor,
		},
	)

	node.Screen.Init(*uiState)

	assert.Equal(t, "golang", string(area.buffer.Buffer()))
	assert.Equal(t, 2, area.caret.Caret())
	assert.Equal(t, 4, area.caret.Anchor())
}

func TestTextArea_Stack(t *testing.T) {
	stack := NewArea().ToNode().Stack

	assert.True(t, stack.Has(NameArea))
}
