package boot

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

	decorator := func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		return func(state state.UIState) {
			called += 1
			next(state)
		}
	}

	mock := screen_test.MockByName("test-node")

	wrapped := Wrap(decorator)(mock)
	assert.Inside(t, Tag, wrapped.Tags)

	wrapped.Screen.Boot(state.UIState{})
	assert.Equal(t, 1, called)
}

func TestWrap_PreservesNextChain(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		return func(state state.UIState) {
			next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Boot: func(state state.UIState) {
			called += 1
		},
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Boot(state.UIState{})

	assert.Equal(t, 1, called)
}

func TestWrap_DoesNotMutateOriginalNode(t *testing.T) {
	tags := set.New[string]()

	mock := screen_test.MockNode{
		Name: "test-node",
		Tags: tags,
	}

	_ = Wrap(func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		return next
	})(mock.ToNode())

	assert.Size(t, 0, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		return func(state state.UIState) {
			captured = target
			next(state)
		}
	}

	mock := screen_test.MockByName("node-123")

	wrapped := Wrap(decorator)(mock)
	wrapped.Screen.Boot(state.UIState{})

	assert.Equal(t, mock.Name, captured.Name)
}

func TestMap_MultipleTransforms(t *testing.T) {
	count := uint(0)

	called0 := uint(0)
	handler1 := func() {
		called0 = count
		count += 1
	}

	called1 := uint(0)
	handler2 := func() {
		called1 = count
		count += 1
	}

	called2 := uint(0)
	handler3 := func() {
		called2 = count
		count += 1
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Boot: func(state.UIState) {},
	}

	node := mock.ToNode()

	node = Map(node, handler1)
	node = Map(node, handler2)
	node = Map(node, handler3)

	node.Screen.Boot(state.UIState{})

	assert.Equal(t, 3, count)

	assert.Equal(t, 0, called0)
	assert.Equal(t, 1, called1)
	assert.Equal(t, 2, called2)
}

func TestUse_ExecutesMiddleware_AndPassesContext(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	name := "test-node"

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.BootFunc]) {
		mwCalled += 1
		assert.Equal(t, name, ctx.Target.Name)
		ctx.Next(uiState)
	}

	mock := screen_test.MockNode{
		Name: name,
		Boot: func(uiState state.UIState) {
			nxCalled += 1
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Boot(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 1, nxCalled)
}

func TestUse_CanShortCircuitChain(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.BootFunc]) {
		mwCalled += 1
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Boot: func(uiState state.UIState) {
			nxCalled += 1
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Boot(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 0, nxCalled)
}
