package checkmenu

import (
	"sort"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap/rw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/checkmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

const Name = "check_menu"

type CheckMenu struct {
	reference    string
	loaded       bool
	bindings     rw.Bindings[CommandRead, CommandWrite]
	definition   rw.Definition
	clock        clock.Clock
	action       *input.CheckAction
	meta         marker.CheckMeta
	distribution style.Distribution
	options      []input.CheckOption
	limit        uint16
	cursor       uint16
}

func New() *CheckMenu {
	return &CheckMenu{
		reference:  Name,
		loaded:     false,
		bindings:   defaultBindings,
		definition: rw.EmptyDefinition(),
		clock:      clock.UnixMilliClock,
		action:     input.EmptyCheckAction(),
		meta:       marker.BracketsCheck,
		options:    make([]input.CheckOption, 0),
		limit:      0,
		cursor:     0,
	}
}

func (n *CheckMenu) Name(name string) *CheckMenu {
	n.reference = name
	return n
}

func (n *CheckMenu) WithWriteBindings(overrides *keymap.Bindings[CommandWrite]) *CheckMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.Write = n.bindings.Write.Overlay(overrides)
	return n
}

func (n *CheckMenu) WithReadBindings(overrides *keymap.Bindings[CommandRead]) *CheckMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.Read = n.bindings.Read.Overlay(overrides)
	return n
}

func (n *CheckMenu) Meta(meta marker.CheckMeta) *CheckMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.meta = meta
	return n
}

func (n *CheckMenu) ActionHandler(handler input.CheckActionHandler) *CheckMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.action.Handler = handler
	return n
}

func (n *CheckMenu) AddOptions(options ...input.CheckOption) *CheckMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	n.options = append(n.options, options...)
	return n
}

func (n *CheckMenu) Cursor(cursor uint16) *CheckMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *CheckMenu) Distribution(distribution style.Distribution) *CheckMenu {
	n.distribution = distribution
	return n
}

func (n *CheckMenu) Limit(limit uint16) *CheckMenu {
	n.limit = limit
	return n
}

func (n *CheckMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *CheckMenu) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loaded = true

	n.loadFromStore(uiState)
	n.definition = rw.DefinitionFromBindings(n.bindings)
}

func (n *CheckMenu) loadFromStore(uiState state.UIState) {
	options, ok := KeyActive.Get(
		uiState.Store,
		n.reference,
	)

	if !ok {
		return
	}

	for i, o := range n.options {
		if options.Has(o.Id) {
			n.switchState(uint16(i), true)
		}
	}

	n.applyLimit()
}

func (n *CheckMenu) keys() screen.Definition {
	return n.definition.Get(n.action.WriteMode)
}

func (n *CheckMenu) tick(uiState *state.UIState, event screen.Event) screen.Result {
	if !n.action.WriteMode {
		return n.tickRead(uiState, event)
	}
	return n.tickWrite(uiState, event)
}

func (n *CheckMenu) tickWrite(uiState *state.UIState, event screen.Event) screen.Result {
	optsLen := uint16(len(n.options))

	switch n.bindings.Write.Command(event.Key.Code) {
	case CmdWriteReadMode:
		n.action.WriteMode = false
	case CmdWriteSwitchState:
		n.switchState(n.cursor)
		n.applyLimit()
		n.tickToStore(uiState)
	case CmdWritePrevOption:
		n.cursor = math.SubClampZero(n.cursor, 1)
	case CmdWriteNextOption:
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = min(optsLen, n.cursor+1)
	case CmdWriteFirstOption:
		n.cursor = 0
	case CmdWriteLastOption:
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = max(0, optsLen)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *CheckMenu) tickToStore(uiState *state.UIState) {
	KeyActive.Set(
		uiState.Store,
		n.reference,
		n.activeIds(),
	)
}

func (n *CheckMenu) tickRead(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.Read.Command(event.Key.Code) {
	case CmdReadWriteMode:
		n.action.WriteMode = true
	}

	return screen.ResultFromUIState(uiState)
}

func (n *CheckMenu) switchState(cursor uint16, state ...bool) *CheckMenu {
	if cursor >= uint16(len(n.options)) {
		return n
	}

	newState := !n.options[cursor].Status
	if len(state) > 0 {
		newState = state[0]
	}

	n.options[cursor].Status = newState

	if n.options[cursor].Status {
		n.options[cursor].Timestamp = n.clock()
	}

	return n
}

func (n *CheckMenu) applyLimit() *CheckMenu {
	if n.limit == 0 {
		return n
	}

	active := make([]*input.CheckOption, 0, len(n.options))
	for i := range n.options {
		if n.options[i].Status {
			active = append(active, &n.options[i])
		}
	}

	if len(active) <= int(n.limit) {
		return n
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].Timestamp < active[j].Timestamp
	})

	excess := len(active) - int(n.limit)
	for i := range excess {
		active[i].Status = false
	}

	return n
}

func (n *CheckMenu) activeIds() set.Set[string] {
	result := set.New[string]()
	for _, v := range n.options {
		if v.Status {
			result.Add(v.Id)
		}
	}
	return result
}

func (n *CheckMenu) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStore(uiState)

	indexmenu := checkmenu.New(n.options).
		WriteMode(n.action.WriteMode).
		Meta(n.meta).
		Cursor(n.cursor)

	vm.Kernel.Push(
		indexmenu.ToUnit(),
	)

	vm.Pager.SetPredicate(
		predicate.Focus(),
	)

	index := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	option := min(index, n.cursor)
	text := n.options[option].Label.Text

	vm.Footer.Push(
		inputline.Wrap(
			drain.UnitFromString(text),
		),
	)

	return *vm
}
