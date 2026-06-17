package screen_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type MockNode struct {
	Name       string
	Definition *screen.Definition
	Keys       screen.KeysFunc
	Boot       screen.BootFunc
	Tick       screen.TickFunc
	View       screen.ViewFunc
	Tags       set.Set[string]
	Stack      set.Set[string]
}

func (t MockNode) ToNode() screen.Node {
	stack := t.Stack
	if t.Stack == nil {
		stack = set.From(t.Name)
	}

	node := screen.NewBuilder().
		Name(t.Name).
		AddStack(stack).
		Boot(
			func(uiState state.UIState) {
				if t.Boot != nil {
					t.Boot(uiState)
				}
			},
		).
		Keys(
			func() screen.KeysFunc {
				if t.Keys != nil {
					return t.Keys
				}

				return func() screen.Definition {
					if t.Definition != nil {
						return *t.Definition
					}

					return screen.EmptyDefinition()
				}
			}()).
		Tick(
			func(s *state.UIState, e screen.Event) screen.Result {
				if t.Tick != nil {
					return t.Tick(s, e)
				}

				return screen.ResultFromUIState(s)
			},
		).
		View(
			func(s state.UIState) viewmodel.ViewModel {
				if t.View != nil {
					return t.View(s)
				}

				return *viewmodel.New()
			},
		).
		ToNode()

	if len(t.Tags) > 0 {
		node.Tags.Merge(t.Tags)
	}

	return node
}

func Helper_ToNode(t *testing.T, node screen.Node) {
	t.Helper()

	assert.NotNil(t, node.Id(), "Node.Stack should be set")
	assert.NotNil(t, node.Name, "Node.Name should be set")
	assert.NotNil(t, node.Stack, "Node.Stack should be set")
	assert.NotNil(t, node.Tags, "Node.Stack should be set")
	assert.NotNil(t, node.Children(), "Node.Stack should be set")

	assert.NotNil(t, node.Screen.Boot, "Screen.Boot should be set")
	assert.NotNil(t, node.Screen.Keys, "Screen.Keys should be set")
	assert.NotNil(t, node.Screen.View, "Screen.View should be set")
	assert.NotNil(t, node.Screen.Tick, "Screen.Tick should be set")
}

func Helper_Propagate(t *testing.T, name string, child uint, node screen.Node) {
	assert.GreaterOrEqual(t, child+1, node.Children())
	assert.True(t, node.Stack.Has(name))
	assert.Equal(t, name, node.Children()[child].Name)
}
