package history

import (
	"fmt"
	"strings"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/trail"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type History struct {
	loaded     bool
	bindings   *keymap.Bindings[Command]
	definition screen.Definition
	trail      *trail.Trail
	meta       marker.HistoryMeta
	node       screen.Node
}

func New(node screen.Node) *History {
	return &History{
		loaded:     false,
		bindings:   defaultBindings,
		definition: screen.EmptyDefinition(),
		trail:      trail.New(trail.DefaultLimit, node),
		meta:       marker.DefaultHistory,
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

func (n *History) SetLimit(limit uint) *History {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.trail = trail.New(limit, n.node)
	return n
}

func (n *History) SetMeta(meta marker.HistoryMeta) *History {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.meta = meta
	return n
}

func (n *History) ToNode() screen.Node {
	snapshot := n.trail.Snapshot().
		ToSlice()

	return screen.NewBuilder().
		Name(n.node.Name).
		AddStack(n.node.Stack).
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		Children(snapshot...).
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
		return n.localTick(uiState, event)
	}

	return n.childTick(uiState, event)
}

func (n *History) localTick(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.Command(event.Key.Code) {
	case CmdPrev:
		back, ok := n.trail.Back()
		if !ok {
			return screen.ResultFromUIState(uiState)
		}

		return screen.ResultFromNode(
			n.makeWrapper(&back),
		)

	case CmdNext:
		next, ok := n.trail.Forward()
		if !ok {
			return screen.ResultFromUIState(uiState)
		}

		return screen.ResultFromNode(
			n.makeWrapper(&next),
		)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *History) childTick(uiState *state.UIState, event screen.Event) screen.Result {
	result := n.node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	n.trail.GoTo(*result.Node)
	result.Node = n.makeWrapper(result.Node)

	return result
}

func (n *History) makeWrapper(node *screen.Node) *screen.Node {
	newWrapper := New(*node)

	newWrapper.loaded = n.loaded
	newWrapper.bindings = n.bindings
	newWrapper.definition = n.definition
	newWrapper.trail = n.trail

	newNode := newWrapper.ToNode()
	return &newNode
}

func (n *History) view(uiState state.UIState) viewmodel.ViewModel {
	vm := n.node.Screen.View(uiState)

	footers := make([]string, 0, 2)

	if back, ok := n.trail.PeekBack(); ok {
		footers = append(footers,
			fmt.Sprintf("%s %s", n.meta.BackTag, back.Name),
		)
	}

	if next, ok := n.trail.PeekForward(); ok {
		footers = append(footers,
			fmt.Sprintf("%s %s", n.meta.NextTag, next.Name),
		)
	}

	if len(footers) == 0 {
		return vm
	}

	line := line.New(
		strings.Join(footers, n.meta.Separator),
		spec.AlignLeft(),
	)

	unit := drain.UnitFromLines(*line).
		AddTag(screen.SystemMetaTag)

	vm.Footer.Unshift(unit)

	return vm
}
