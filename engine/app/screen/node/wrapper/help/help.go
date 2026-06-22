package help

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/help"
)

type Help struct {
	bindings *keymap.Bindings[Command]
	visible  bool
	node     screen.Node
}

func New(node screen.Node) *Help {
	return &Help{
		bindings: defaultBindings,
		visible:  false,
		node:     node,
	}
}

func (n *Help) WithBindings(overrides *keymap.Bindings[Command]) *Help {
	n.bindings = n.bindings.Overlay(overrides)
	return n
}

func (n *Help) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Boot(n.node.Screen.Boot).
		Keys(n.node.Screen.Keys).
		Tick(n.tick).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *Help) tick(uiState *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Keys()
	if !definition.IsRequired(event.Key) {
		return n.localTick(uiState, event)
	}

	return n.childTick(uiState, event)
}

func (n *Help) localTick(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.Command(event.Key.Code) {
	case CmdSwitchDisplay:
		n.visible = !n.visible
	}

	uiState.Helper.ShowHelp = n.visible
	return screen.ResultFromUIState(uiState)
}

func (n *Help) childTick(uiState *state.UIState, event screen.Event) screen.Result {
	n.visible = uiState.Helper.ShowHelp

	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.visible = n.visible
	newScreen := newWrapper.ToNode()
	result.Node = &newScreen

	return result
}

func (n *Help) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)
	if !n.visible {
		return vm
	}

	definition := n.node.Screen.Keys()

	vm.Footer.Push(
		help.UnitFromFields(definition.Descriptor.ToValuesSlice()),
	)

	return vm
}
