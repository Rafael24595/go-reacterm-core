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
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestWrap_AddsTag_AndExecutesDecorator(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(state state.UIState) viewmodel.ViewModel {
			called += 1
			return next(state)
		}
	}

	mock := screen_test.MockByName("test-node")

	wrapped := Wrap(decorator)(mock)
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

	assert.Empty(t, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(state state.UIState) viewmodel.ViewModel {
			captured = target
			return next(state)
		}
	}

	mock := screen_test.MockByName("node-123")

	wrapped := Wrap(decorator)(mock)
	wrapped.Screen.View(state.UIState{})

	assert.Equal(t, mock.Name, captured.Name)
}

func TestMap_Transforms(t *testing.T) {
	mockUnit := drawable_test.MockUnit{}

	handler := func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		vm.Kernel.Push(
			mockUnit.ToUnit(),
		)
		return vm
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		View: func(uiState state.UIState) viewmodel.ViewModel {
			return *viewmodel.New()
		},
	}

	node := Map(mock.ToNode(), handler)

	vm := node.Screen.View(state.UIState{})

	assert.Equal(t, 1, vm.Kernel.Size())
}

func TestMap_MultipleTransforms(t *testing.T) {
	count := uint(0)

	called0 := uint(0)
	handler1 := func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		called0 = count
		count += 1
		return vm
	}

	called1 := uint(0)
	handler2 := func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		called1 = count
		count += 1
		return vm
	}

	called2 := uint(0)
	handler3 := func(vm viewmodel.ViewModel) viewmodel.ViewModel {
		called2 = count
		count += 1
		return vm
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		View: func(uiState state.UIState) viewmodel.ViewModel {
			return *viewmodel.New()
		},
	}

	node := mock.ToNode()

	node = Map(node, handler1)
	node = Map(node, handler2)
	node = Map(node, handler3)

	node.Screen.View(state.UIState{})

	assert.Equal(t, 3, count)

	assert.Equal(t, 0, called0)
	assert.Equal(t, 1, called1)
	assert.Equal(t, 2, called2)
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
