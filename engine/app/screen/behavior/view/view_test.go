package view

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestWrap_AddsTag_AndExecutesDecorator(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(state state.UIState) viewmodel.ViewModel {
			called += 1
			return next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	assert.Inside(t, Tag, wrapped.Tags)

	wrapped.Screen.View(state.UIState{})
	assert.Equal(t, 1, called)
}

func TestWrap_PreservesNextChain(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(state state.UIState) viewmodel.ViewModel {
			return next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		View: func(state state.UIState) viewmodel.ViewModel {
			called += 1
			return *viewmodel.New()
		},
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.View(state.UIState{})

	assert.Equal(t, 1, called)
}

func TestWrap_DoesNotMutateOriginalNode(t *testing.T) {
	tags := set.New[string]()

	mock := screen_test.MockNode{
		Name: "test-node",
		Tags: tags,
	}

	_ = Wrap(func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return next
	})(mock.ToNode())

	assert.Size(t, 0, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(state state.UIState) viewmodel.ViewModel {
			captured = target
			return next(state)
		}
	}

	mock := screen_test.MockNode{
		Name: "node-123",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.View(state.UIState{})

	assert.Equal(t, mock.Name, captured.Name)
}

func TestUse_ExecutesMiddleware_AndPassesContext(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	name := "test-node"

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.ViewFunc]) viewmodel.ViewModel {
		mwCalled += 1
		assert.Equal(t, name, ctx.Target.Name)
		return ctx.Next(uiState)
	}

	mock := screen_test.MockNode{
		Name: name,
		View: func(uiState state.UIState) viewmodel.ViewModel {
			nxCalled += 1
			return viewmodel.ViewModel{}
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.View(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 1, nxCalled)
}

func TestUse_CanShortCircuitChain(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(uiState state.UIState, ctx behavior.Context[screen.ViewFunc]) viewmodel.ViewModel {
		mwCalled += 1
		return viewmodel.ViewModel{}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		View: func(uiState state.UIState) viewmodel.ViewModel {
			nxCalled += 1
			return viewmodel.ViewModel{}
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.View(state.UIState{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 0, nxCalled)
}
