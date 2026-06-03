package tick

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

	decorator := func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		return func(state *state.UIState, event screen.Event) screen.Result {
			called += 1
			return next(state, event)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	assert.Contains(t, wrapped.Tags, Tag)

	wrapped.Screen.Tick(&state.UIState{}, screen.Event{})
	assert.Equal(t, 1, called)
}

func TestWrap_PreservesNextChain(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		return func(state *state.UIState, event screen.Event) screen.Result {
			return next(state, event)
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Tick: func(state *state.UIState, event screen.Event) screen.Result {
			called += 1
			return screen.EmptyResult()
		},
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, 1, called)
}

func TestWrap_DoesNotMutateOriginalNode(t *testing.T) {
	tags := set.New[string]()

	mock := screen_test.MockNode{
		Name: "test-node",
		Tags: tags,
	}

	_ = Wrap(func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		return next
	})(mock.ToNode())

	assert.Len(t, 0, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		return func(state *state.UIState, event screen.Event) screen.Result {
			captured = target
			return next(state, event)
		}
	}

	mock := screen_test.MockNode{
		Name: "node-123",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, mock.Name, captured.Name)
}
