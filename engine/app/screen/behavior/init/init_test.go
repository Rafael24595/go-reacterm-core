package init

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestWrap_AddsTag_AndExecutesDecorator(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.InitFunc) screen.InitFunc {
		return func(state state.UIState) {
			called += 1
			next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	assert.Inside(t, Tag, wrapped.Tags)

	wrapped.Screen.Init(state.UIState{})
	assert.Equal(t, 1, called)
}

func TestWrap_PreservesNextChain(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.InitFunc) screen.InitFunc {
		return func(state state.UIState) {
			next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Init: func(state state.UIState) {
			called += 1
		},
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Init(state.UIState{})

	assert.Equal(t, 1, called)
}

func TestWrap_DoesNotMutateOriginalNode(t *testing.T) {
	tags := set.New[string]()

	mock := screen_test.MockNode{
		Name: "test-node",
		Tags: tags,
	}

	_ = Wrap(func(target behavior.Target, next screen.InitFunc) screen.InitFunc {
		return next
	})(mock.ToNode())

	assert.Size(t, 0, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.InitFunc) screen.InitFunc {
		return func(state state.UIState) {
			captured = target
			next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "node-123",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Init(state.UIState{})

	assert.Equal(t, mock.Name, captured.Name)
}

func TestUse_ExecutesMiddleware_AndPassesContext(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	name := "test-node"

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.InitFunc])  {
		mwCalled += 1
		assert.Equal(t, name, ctx.Target.Name)
		ctx.Next(uiState)
	}

	mock := screen_test.MockNode{
		Name: name,
		Init: func(uiState state.UIState) {
			nxCalled += 1
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Init(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 1, nxCalled)
}

func TestUse_CanShortCircuitChain(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.InitFunc]) {
		mwCalled += 1
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Init: func(uiState state.UIState) {
			nxCalled += 1
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Init(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 0, nxCalled)
}
