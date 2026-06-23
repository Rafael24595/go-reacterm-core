package history

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type History struct {
	loaded     bool
	bindings   *keymap.Bindings[Command]
	definition screen.Definition
	history    *screen.Node
	node       screen.Node
}

func New(node screen.Node) *History {
	return &History{
		loaded:     false,
		bindings:   defaultBindings,
		definition: screen.EmptyDefinition(),
		node:       node,
	}
}

func (n *History) WithBindings(overrides *keymap.Bindings[Command]) *History {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings = n.bindings.Overlay(overrides)
	return n
}

func (n *History) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *History) boot(uiState state.UIState) {
	if !n.loaded {
		n.loaded = true
		n.definition = keymap.BindingsToDefinition(n.bindings)
	}

	n.node.Screen.Boot(uiState)
}

func (n *History) keys() screen.Definition {
	return n.definition.Merge(
		n.node.Screen.Keys(),
	)
}

func (n *History) tick(uiState *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Keys()
	if !definition.IsRequired(event.Key) {
		result := n.localTick(uiState, event)
		if result != nil {
			return *result
		}
	}

	return n.childTick(uiState, event)
}

func (n *History) localTick(_ *state.UIState, event screen.Event) *screen.Result {
	command := n.bindings.Command(event.Key.Code)
	if n.history == nil || command != CmdBack {
		return nil
	}

	result := screen.ResultFromNode(
		n.makeWrapper(n.history),
	)

	return &result
}

func (n *History) childTick(uiState *state.UIState, event screen.Event) screen.Result {
	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	result.Node = n.makeWrapper(result.Node)
	return result
}

func (n *History) makeWrapper(node *screen.Node) *screen.Node {
	newWrapper := New(*node)

	newWrapper.loaded = n.loaded
	newWrapper.bindings = n.bindings
	newWrapper.definition = n.definition
	newWrapper.history = &n.node

	newNode := newWrapper.ToNode()
	return &newNode
}

func (n *History) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)

	if n.history == nil {
		return vm
	}

	footer := text.NewLine(
		fmt.Sprintf("back: %s", n.history.Name),
		spec.AlignLeft(),
	)

	vm.Footer.Unshift(
		drain.UnitFromLines(*footer).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}
