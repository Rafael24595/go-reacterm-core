package pagination

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const errf_unhandled = "unhandled pager type '%d'"

type Pagination struct {
	actionKind  action.Kind
	forceAction *action.Action
	node        screen.Node
}

func New(screen screen.Node) *Pagination {
	return &Pagination{
		actionKind:  action.KindPaged,
		forceAction: nil,
		node:        screen,
	}
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
		Boot(n.node.Screen.Boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		Children(n.node).
		ToNode()
}

func (n *Pagination) keys() screen.Definition {
	node := n.node.Screen.Keys()
	return base_definition.Merge(
		n.findDefinition().Merge(node),
	)
}

func (n *Pagination) findDefinition() screen.Definition {
	if source, ok := definitions[n.actionKind]; ok {
		return source
	}

	assert.Unreachable("unhandled action definition %d", n.actionKind)
	return pager_definition
}

func (n *Pagination) tick(uiState *state.UIState, event screen.Event) screen.Result {
	definition := n.node.Screen.Keys()

	if !definition.IsRequired(event.Key) {
		result := n.localTick(uiState, event)
		if result != nil {
			return *result
		}
	}

	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	newWrapper := New(*result.Node)
	newWrapper.actionKind = n.actionKind
	newWrapper.forceAction = n.forceAction
	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Pagination) localTick(uiState *state.UIState, event screen.Event) *screen.Result {
	keys, ok := keys[n.actionKind]

	assert.True(ok, errf_unhandled, action.KindPaged)

	if event.Key.Code == key.ActionPageUp || event.Key.Code == keys.back {
		uiState.Pager.DecTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	}

	if event.Key.Code == key.ActionPageDown || event.Key.Code == keys.next {
		uiState.Pager.IncTarget()
		result := screen.ResultFromUIState(uiState)
		return &result
	}

	return nil
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

	footer := []text.Line{
		*text.NewLine(
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
