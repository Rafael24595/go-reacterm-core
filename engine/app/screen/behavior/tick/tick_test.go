package tick

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

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
	assert.Inside(t, Tag, wrapped.Tags)

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

	assert.Size(t, 0, tags)
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

func TestMap_Transforms(t *testing.T) {
	handler := func(result screen.Result) screen.Result {
		mock := screen_test.MockNode{}
		node := mock.ToNode()

		result.Node = &node

		return result
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Tick: func(*state.UIState, screen.Event) screen.Result {
			return screen.EmptyResult()
		},
	}

	node := Map(mock.ToNode(), handler)

	result := node.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.NotNil(t, 1, result.Node)
}

func TestMap_MultipleTransforms(t *testing.T) {
	count := uint(0)

	called0 := uint(0)
	handler1 := func(result screen.Result) screen.Result {
		called0 = count
		count += 1
		return result
	}

	called1 := uint(0)
	handler2 := func(result screen.Result) screen.Result {
		called1 = count
		count += 1
		return result
	}

	called2 := uint(0)
	handler3 := func(result screen.Result) screen.Result {
		called2 = count
		count += 1
		return result
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

	node.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, 3, count)

	assert.Equal(t, 0, called0)
	assert.Equal(t, 1, called1)
	assert.Equal(t, 2, called2)
}

func TestUse_ExecutesMiddleware_AndPassesContext(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	name := "test-node"

	middleware := func(uiState *state.UIState, event screen.Event, ctx behavior.Context[screen.TickFunc]) screen.Result {
		mwCalled += 1
		assert.Equal(t, name, ctx.Target.Name)
		return ctx.Next(uiState, event)
	}

	mock := screen_test.MockNode{
		Name: name,
		Tick: func(uiState *state.UIState, event screen.Event) screen.Result {
			nxCalled += 1
			return screen.EmptyResult()
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 1, nxCalled)
}

func TestUse_CanShortCircuitChain(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(uiState *state.UIState, event screen.Event, ctx behavior.Context[screen.TickFunc]) screen.Result {
		mwCalled += 1
		return screen.EmptyResult()
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Tick: func(uiState *state.UIState, event screen.Event) screen.Result {
			nxCalled += 1
			return screen.EmptyResult()
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Tick(&state.UIState{}, screen.Event{})

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 0, nxCalled)
}

func TestOnKey_ExecutesMiddleware_WhenKeyMatches(t *testing.T) {
	called := uint(0)

	middleware := func(uiState *state.UIState, event screen.Event, ctx behavior.Context[screen.TickFunc]) screen.Result {
		called += 1
		return screen.EmptyResult()
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	node := OnKey(mock.ToNode(), middleware, key.ActionEnter, key.ActionEsc)

	matchingEvent := screen.Event{
		Key: key.Key{
			Code: key.ActionEnter,
		},
	}

	node.Screen.Tick(&state.UIState{}, matchingEvent)

	assert.Equal(t, 1, called)
}

func TestOnKey_BypassesMiddleware_AndCallsNext_WhenKeyDoesNotMatch(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(uiState *state.UIState, event screen.Event, ctx behavior.Context[screen.TickFunc]) screen.Result {
		mwCalled += 1
		return screen.EmptyResult()
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Tick: func(uiState *state.UIState, event screen.Event) screen.Result {
			nxCalled += 1
			return screen.EmptyResult()
		},
	}

	node := OnKey(mock.ToNode(), middleware, key.ActionEnter, key.ActionEsc)

	matchingEvent := screen.Event{
		Key: key.Key{
			Code: key.ActionDelete,
		},
	}

	node.Screen.Tick(&state.UIState{}, matchingEvent)

	assert.Equal(t, 0, mwCalled)
	assert.Equal(t, 1, nxCalled)
}
