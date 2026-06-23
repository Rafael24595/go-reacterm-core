package form

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/dummy"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/gutter"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/form"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "form"

type Form struct {
	reference  string
	loaded     bool
	bindings   bindings
	definition definition
	pointer    uint8
	focused    bool
	cursor     uint16
	items      []entry.Entry
	dummies    set.Set[int]
}

func New() *Form {
	return &Form{
		reference:  Name,
		loaded:     false,
		bindings:   defaultBindings,
		definition: emptyDefinition(),
		pointer:    0,
		focused:    false,
		cursor:     0,
		items:      make([]entry.Entry, 0),
		dummies:    set.New[int](),
	}
}

func (n *Form) AddNode(node screen.Node, opts ...entry.Option) *Form {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	n.items = append(n.items,
		entry.New(node, opts...),
	)
	return n
}

func (n *Form) AddBreak(rows ...winsize.Rows) *Form {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	fixed := winsize.Rows(1)
	if len(rows) > 0 {
		fixed = rows[0]
	}

	n.dummies.Add(len(n.items))

	return n.AddNode(
		dummy.ToNode(),
		entry.WithLayout(
			layer.Fixed(fixed),
		),
	)
}

func (n *Form) ToNode() screen.Node {
	n.setCursor(0)
	n.cursor = n.incCursor(0)

	builder := screen.NewBuilder().
		Name(n.reference).
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view)

	for _, v := range n.items {
		builder.Children(v.Node).
			AddStack(v.Node.Stack)
	}

	return builder.ToNode()
}

func (n *Form) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loaded = true
	n.definition = definitionFromBindings(n.bindings)

	for _, item := range n.items {
		item.Node.Screen.Boot(uiState)
	}
}

func (n *Form) keys() screen.Definition {
	local := n.definition.get(n.focused)
	if !n.focused {
		return local
	}

	focus, ok := n.focusItem()
	if !ok || !focus.Selectable {
		return local
	}

	return local.Merge(
		focus.Node.Screen.Keys(),
	)
}

func (n *Form) tick(uiState *state.UIState, event screen.Event) screen.Result {
	if !n.focused {
		return n.readTick(uiState, event)
	}
	return n.writeTick(uiState, event)
}

func (n *Form) writeTick(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.write.Command(event.Key.Code) {
	case CmdWriteReadMode:
		n.focused = false
	}

	return n.tryFocusTick(uiState, event)
}

func (n *Form) readTick(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.read.Command(event.Key.Code) {
	case CmdReadWriteMode:
		n.focused = true
		return n.tryFocusTick(uiState, event)
	case CmdReadPrevOption:
		n.cursor = n.decCursor()
	case CmdReadNextOption:
		n.cursor = n.incCursor()
	case CmdReadFirstOption:
		n.setCursor(0)
		n.cursor = n.incCursor(0)
	case CmdReadLastOption:
		items := len(n.items)
		n.setCursor(uint16(items))
		n.cursor = n.decCursor(0)
	case CmdReadSwitchPointer:
		n.pointer = form.NextPointer(n.pointer)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *Form) tryFocusTick(uiState *state.UIState, event screen.Event) screen.Result {
	focus, ok := n.focusItem()
	if !ok {
		return screen.ResultFromUIState(uiState)
	}

	definition := focus.Node.Screen.Keys()
	if !definition.IsRequired(event.Key) {
		return screen.ResultFromUIState(uiState)
	}

	return n.focusTick(uiState, event, focus)
}

func (n *Form) focusTick(uiState *state.UIState, event screen.Event, focus entry.Entry) screen.Result {
	result := focus.Node.Screen.Tick(uiState, event)
	if result.Node == nil {
		return result
	}

	newItems := make([]entry.Entry, len(n.items))
	copy(newItems, n.items)

	newWrapper := New()

	newWrapper.reference = n.reference
	newWrapper.pointer = n.pointer
	newWrapper.focused = n.focused
	newWrapper.cursor = n.cursor
	newWrapper.items = newItems

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
}

func (n *Form) setCursor(cursor uint16) uint16 {
	items := len(n.items)
	if items == 0 {
		return 0
	}

	n.cursor = min(cursor, uint16(items-1))
	return n.cursor
}

func (n *Form) incCursor(data ...uint16) uint16 {
	inc := 1
	if len(data) > 0 {
		inc = int(data[0])
	}

	return n.moveCursor(true, inc)
}

func (n *Form) decCursor(data ...uint16) uint16 {
	dec := 1
	if len(data) > 0 {
		dec = int(data[0])
	}

	return n.moveCursor(false, -dec)
}

func (n *Form) moveCursor(sign bool, step int) uint16 {
	size := len(n.items)
	if size == 0 {
		return 0
	}

	base := int(n.cursor)

	dir := 1
	if !sign {
		dir = -1
	}

	for i := range size {
		cursor := (base + step + dir*i) % size
		if cursor < 0 {
			cursor += size
		}

		if n.items[cursor].Selectable {
			return uint16(cursor)
		}
	}

	return n.cursor
}

func (n *Form) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	pointer := form.FindPointer(n.pointer)

	// TODO: Compile headers and footers?
	for i, e := range n.items {
		cvm := e.Node.Screen.View(uiState)

		opts := make([]gutter.Option, 0, 1)

		if pointer.HasNone(form.PointerGutter) || n.cursor != uint16(i) {
			opts = append(opts,
				gutter.WithLeftGutter(gutter.DefaultEmpty),
			)
		}

		unit := gutter.Unit(
			cvm.Kernel.ToUnit(),
			opts...,
		)

		layer := layer.New(unit, e.Opts...)
		vm.Kernel.PushLayer(layer)

		if cvm.Behavior.NeedsPulse {
			vm.Behavior.NeedsPulse = true
		}

		// TODO: Improve pager inheritance policy.
		// TODO: Inherit headers?
		if n.cursor == uint16(i) {
			vm.Footer.Push(cvm.Footer.Units()...)

			if !layer.Static() && !n.dummies.Has(i) {
				vm.Pager = cvm.Pager
			}
		}
	}

	focus, ok := n.focusItem()
	if ok && pointer.HasAny(form.PointerPrompt) {
		label := text.NewFragment(focus.Node.Name).
			AddAtom(atom.Select)

		vm.Footer.Push(
			inputline.FromFragment(*label),
		)
	}

	return *vm
}

func (n *Form) focusItem() (entry.Entry, bool) {
	if n.cursor >= uint16(len(n.items)) {
		return entry.Entry{}, false
	}

	return n.items[n.cursor], true
}
