package keys

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestWrap_AddsTag_AndExecutesDecorator(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		return func() screen.Definition {
			called += 1
			return next()
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	assert.Inside(t, Tag, wrapped.Tags)

	wrapped.Screen.Keys()
	assert.Equal(t, 1, called)
}

func TestWrap_PreservesNextChain(t *testing.T) {
	called := uint(0)

	decorator := func(target behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		return func() screen.Definition {
			return next()
		}
	}

	mock := screen_test.MockNode{
		Name: "test-node",
	}

	node := mock.ToNode()
	node.Screen.Keys = func() screen.Definition {
		called += 1
		return screen.DefinitionFromActions()
	}

	wrapped := Wrap(decorator)(node)
	wrapped.Screen.Keys()

	assert.Equal(t, 1, called)
}

func TestWrap_DoesNotMutateOriginalNode(t *testing.T) {
	tags := set.New[string]()

	mock := screen_test.MockNode{
		Name: "test-node",
		Tags: tags,
	}

	_ = Wrap(func(target behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		return next
	})(mock.ToNode())

	assert.Size(t, 0, tags)
}

func TestWrap_TargetIsCorrect(t *testing.T) {
	captured := behavior.Target{}

	decorator := func(target behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		return func() screen.Definition {
			captured = target
			return next()
		}
	}

	mock := screen_test.MockNode{
		Name: "node-123",
	}

	wrapped := Wrap(decorator)(mock.ToNode())
	wrapped.Screen.Keys()

	assert.Equal(t, mock.Name, captured.Name)
}

func TestMap_Transforms(t *testing.T) {
	handler := func(definition screen.Definition) screen.Definition {
		return screen.DefinitionFromActions(
			key.ActionAll,
		)
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Keys: func() screen.Definition {
			return screen.EmptyDefinition()
		},
	}

	node := Map(mock.ToNode(), handler)

	definition := node.Screen.Keys()

	assert.Equal(t, 1, definition.RequireKeys.Size())
}

func TestMap_MultipleTransforms(t *testing.T) {
	count := uint(0)

	called0 := uint(0)
	handler1 := func(definition screen.Definition) screen.Definition {
		called0 = count
		count += 1
		return definition
	}

	called1 := uint(0)
	handler2 := func(definition screen.Definition) screen.Definition {
		called1 = count
		count += 1
		return definition
	}

	called2 := uint(0)
	handler3 := func(definition screen.Definition) screen.Definition {
		called2 = count
		count += 1
		return definition
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Keys: func() screen.Definition {
			return screen.EmptyDefinition()
		},
	}

	node := mock.ToNode()

	node = Map(node, handler1)
	node = Map(node, handler2)
	node = Map(node, handler3)

	node.Screen.Keys()

	assert.Equal(t, 3, count)

	assert.Equal(t, 0, called0)
	assert.Equal(t, 1, called1)
	assert.Equal(t, 2, called2)
}

func TestUse_ExecutesMiddleware_AndPassesContext(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	name := "test-node"

	middleware := func(ctx behavior.Context[screen.KeysFunc]) screen.Definition {
		mwCalled += 1
		assert.Equal(t, name, ctx.Target.Name)
		return ctx.Next()
	}

	mock := screen_test.MockNode{
		Name: name,
		Keys: func() screen.Definition {
			nxCalled += 1
			return screen.EmptyDefinition()
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Keys()

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 1, nxCalled)
}

func TestUse_CanShortCircuitChain(t *testing.T) {
	mwCalled := uint(0)
	nxCalled := uint(0)

	middleware := func(ctx behavior.Context[screen.KeysFunc]) screen.Definition {
		mwCalled += 1
		return screen.EmptyDefinition()
	}

	mock := screen_test.MockNode{
		Name: "test-node",
		Keys: func() screen.Definition {
			nxCalled += 1
			return screen.EmptyDefinition()
		},
	}

	node := Use(mock.ToNode(), middleware)
	node.Screen.Keys()

	assert.Equal(t, 1, mwCalled)
	assert.Equal(t, 0, nxCalled)
}
