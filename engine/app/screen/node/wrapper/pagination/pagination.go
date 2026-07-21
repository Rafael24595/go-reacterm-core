package pagination

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const errf_unhandled = "unhandled pager type '%d'"

type Pagination struct {
	loaded      bool
	bindings    bindings
	definition  definition
	actionKind  action.Kind
	forceAction *action.Action
	node        screen.Node
}

func New(screen screen.Node) *Pagination {
	return &Pagination{
		loaded:      false,
		bindings:    defaultBindings,
		definition:  emptyDefinition(),
		actionKind:  action.KindPaged,
		forceAction: nil,
		node:        screen,
	}
}

func (n *Pagination) WithBaseBindings(overrides *keymap.Bindings[Command]) *Pagination {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.base = n.bindings.base.Overlay(overrides)
	return n
}

func (n *Pagination) WithBindingsForPaged(overrides *keymap.Bindings[Command]) *Pagination {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.pager = n.bindings.pager.Overlay(overrides)
	return n
}

func (n *Pagination) WithBindingsForScroll(overrides *keymap.Bindings[Command]) *Pagination {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.scroll = n.bindings.scroll.Overlay(overrides)
	return n
}

func (n *Pagination) ForceEngine(forceAction action.Action) *Pagination {
	n.forceAction = &forceAction
	n.actionKind = forceAction.Kind

	return n
}

func (n *Pagination) ToNode() screen.Node {
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

func (n *Pagination) boot(uiState state.UIState) {
	if !n.loaded {
		n.loaded = true
		n.definition = definitionFromBindings(n.bindings)
	}

	n.node.Screen.Boot(uiState)
}

func (n *Pagination) keys() screen.Definition {
	return n.definition.get(n.actionKind).
		Merge(n.node.Screen.Keys())
}

func (n *Pagination) tick(uiState *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Keys()

	if !definition.IsRequired(event.Key) {
		result := n.localTick(uiState, event)
		if result != nil {
			return *result
		}
	}

	return n.childTick(uiState, event)
}

func (n *Pagination) localTick(uiState *state.UIState, event screen.Event) *screen.Result {
	binding := n.bindings.get(n.actionKind)

	switch binding.Command(event.Key.Code) {
	case CmdPageUp, CmdPrevPage:
		uiState.Pager.DecTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	case CmdPageDown, CmdNextPage:
		uiState.Pager.IncTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	}

	return nil
}

func (n *Pagination) childTick(uiState *state.UIState, event screen.Event) screen.Result {
	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)

	newWrapper.loaded = n.loaded
	newWrapper.bindings = n.bindings
	newWrapper.definition = n.definition
	newWrapper.actionKind = n.actionKind
	newWrapper.forceAction = n.forceAction

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Pagination) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)
	if n.forceAction != nil {
		vm.Pager.SetAction(*n.forceAction)
	}

	n.actionKind = vm.Pager.Action.Kind

	if !n.shouldShowPage(uiState, vm) {
		return vm
	}

	label, ok := labels[n.actionKind]

	assert.True(ok, errf_unhandled, n.actionKind)

	footer := []line.Line{
		line.TextSpec(
			fmt.Sprintf("%s: %d", label, uiState.Pager.ActualPage),
			spec.AlignLeft(),
		),
	}

	vm.Footer.Unshift(
		drain.UnitFromLines(footer...).
			AddTag(screen.SystemMetaTag),
	)

	return vm
}

func (n *Pagination) shouldShowPage(uiState state.UIState, vm viewmodel.ViewModel) bool {
	if vm.Pager.Predicate.Kind != predicate.KindPage {
		return false
	}

	if uiState.Pager.ForceShow {
		return true
	}

	return uiState.Pager.HasMore || uiState.Pager.ActualPage > 0
}
