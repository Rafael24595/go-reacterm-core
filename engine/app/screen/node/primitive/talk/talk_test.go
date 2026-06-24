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

func TestTalk_Boot(t *testing.T) {
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

	talk := New().AddMessage(messges...)
	node := talk.ToNode()

	assert.Equal(t, 0, talk.cursor)
	assert.Size(t, 4, talk.messages)

	uiState := state.NewUIState()

	cursor := uint16(4)
	newMessages := append(messges,
		chat.Message{
			Message: "message_05",
		})

	KeySync.Set(
		uiState.Store,
		node.Name,
		Sync{
			Cursor:   &cursor,
			Messages: &newMessages,
		},
	)

	node.Screen.Boot(*uiState)

	_, ok := KeySync.Get(uiState.Store, talk.reference)

	assert.False(t, ok)

	assert.Equal(t, 4, talk.cursor)
	assert.Size(t, 5, talk.messages)
}

func TestTalk_Stack(t *testing.T) {
	stack := New().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
