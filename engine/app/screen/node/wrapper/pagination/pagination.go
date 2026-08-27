package pagination

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
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
	loaded     bool
	bindings   bindings
	definition definition
	stepKind   step.Kind
	forceStep  *step.Step
	node       screen.Node
}

func New(node screen.Node) *Pagination {
	return &Pagination{
		loaded:     false,
		bindings:   defaultBindings,
		definition: emptyDefinition(),
		stepKind:   step.KindPage,
		forceStep:  nil,
		node:       node,
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

func (n *Pagination) WithStep(step step.Step) *Pagination {
	n.forceStep = &step
	n.stepKind = step.Kind

	return n
}

func (n *Pagination) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.node.Name).
		Stack(n.node.Stack).
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
	return n.definition.get(n.stepKind).
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
	binding := n.bindings.get(n.stepKind)

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

	node, hasNode := result.TryGetNode()
	if !hasNode {
		return result
	}

	newWrapper := New(node)

	newWrapper.loaded = n.loaded
	newWrapper.bindings = n.bindings
	newWrapper.definition = n.definition
	newWrapper.stepKind = n.stepKind
	newWrapper.forceStep = n.forceStep

	newNode := newWrapper.ToNode()
	result.SetNode(newNode)

	return result
}

func (n *Pagination) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)
	if n.forceStep != nil {
		vm.Pager.WithStep(*n.forceStep)
	}

	n.stepKind = vm.Pager.Step.Kind

	if !n.shouldShowPage(uiState, vm) {
		return vm
	}

	label, ok := labels[n.stepKind]

	assert.True(ok, errf_unhandled, n.stepKind)

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
	if vm.Pager.Rule.Kind != rule.KindPage {
		return false
	}

	if uiState.Pager.ForceShow {
		return true
	}

	return uiState.Pager.HasMore || uiState.Pager.ActualPage > 0
}
