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
	assert.Contains(t, wrapped.Tags, Tag)

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

	assert.Len(t, 0, tags)
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
