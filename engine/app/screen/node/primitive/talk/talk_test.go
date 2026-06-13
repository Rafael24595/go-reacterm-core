package talk

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTalk_ToNode(t *testing.T) {
	node := New().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)

	assert.Equal(t, node.Name, "base")
}

func TestTalk_Init(t *testing.T) {
	messges := []chat.Message{
		{
			Message: "message_03",
		},
		{
			Message: "message_01",
		},
		{
			Message: "message_03",
		},
		{
			Message: "message_04",
		},
	}

	menu := New().AddMessage(messges...)
	node := menu.ToNode()

	assert.Equal(t, 0, menu.cursor)
	assert.Size(t, 4, menu.messages)

	uiState := state.NewUIState()

	KeyCursor.Set(
		uiState.Store,
		node.Name,
		4,
	)

	KeyMessages.Set(
		uiState.Store,
		node.Name,
		append(messges,
			chat.Message{
				Message: "message_05",
			},
		),
	)

	node.Screen.Init(*uiState)

	assert.Equal(t, 4, menu.cursor)
	assert.Size(t, 5, menu.messages)
}

func TestTalk_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
