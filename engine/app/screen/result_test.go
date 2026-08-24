package screen

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestEmptyResult(t *testing.T) {
	res := EmptyResult()

	assert.False(t, res.HasNode())
	assert.True(t, IsZeroNode(res.GetNode()))

	node, ok := res.TryGetNode()
	assert.False(t, ok)
	assert.True(t, IsZeroNode(node))
}

func TestResultFromNode(t *testing.T) {
	mock := NewBuilder().
		NameAsStack("mock_node").
		ToNode()

	res := ResultFromNode(mock)

	assert.True(t, res.HasNode())
	assert.Equal(t, mock.id, res.GetNode().id)

	node, ok := res.TryGetNode()
	assert.True(t, ok)
	assert.Equal(t, mock.id, node.id)
}

func TestResultFromUIState(t *testing.T) {
	t.Run("Valid UIState pointer", func(t *testing.T) {
		uiState := &state.UIState{
			Pager: state.PagerContext{},
		}

		res := ResultFromUIState(uiState)

		assert.False(t, res.HasNode())
		assert.Equal(t, uiState.Pager, res.Pager)
	})

	t.Run("Nil UIState pointer handles gracefully", func(t *testing.T) {
		res := ResultFromUIState(nil)

		assert.False(t, res.HasNode())
		assert.Equal(t, state.PagerContext{}, res.Pager)
	})
}

func TestResult_SetNode(t *testing.T) {
	mock := EmptyResult()
	
	targetMock := NewBuilder().
		NameAsStack("mock_target").
		ToNode()

	chainRes := mock.SetNode(targetMock)

	assert.True(t, mock.HasNode())
	assert.Equal(t, targetMock.id, mock.GetNode().id)
	assert.Equal(t, &mock, chainRes)
}
