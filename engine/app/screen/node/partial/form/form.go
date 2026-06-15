package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/dummy"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/entry"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/gutter"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/form"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "form"

type Form struct {
	reference string
	pointer   uint8
	focused   bool
	cursor    uint16
	steps     []pipeline.Transformer
	items     []entry.Entry
}

func New() *Form {
	return &Form{
		reference: Name,
		pointer:   0,
		focused:   false,
		cursor:    0,
		steps:     make([]pipeline.Transformer, 0),
		items:     make([]entry.Entry, 0),
	}
}

func (n *Form) PushSteps(steps ...pipeline.Transformer) *Form {
	n.steps = append(n.steps, steps...)
	return n
}

func (n *Form) AddNode(
	node screen.Node,
	opts ...entry.Option,
) *Form {
	n.items = append(n.items,
		entry.New(node, opts...),
	)
	return n
}

func (n *Form) AddBreak(rows ...winsize.Rows) *Form {
	fixed := winsize.Rows(1)
	if len(rows) > 0 {
		fixed = rows[0]
	}

	return n.AddNode(
		dummy.ToNode(),
		entry.WithLayout(
			layer.Fixed(fixed),
		),
	)
}

func (n *Form) ToNode() screen.Node {
	builder := screen.NewBuilder().
		Name(n.reference).
		Init(n.init).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view)

	for _, v := range n.items {
		builder.Children(v.Node).
			AddStack(v.Node.Stack)
	}

	return builder.ToNode()
}

func (n *Form) init(uiState state.UIState) {
	for _, item := range n.items {
		item.Node.Screen.Init(uiState)
	}
}

func (n *Form) keys() screen.Definition {
	local := sources

	item := n.items[n.cursor]
	if item.Selectable {
		local = local.Merge(
			item.Node.Screen.Keys(),
		)
	}

	return local
}

func (n *Form) tick(uiState *state.UIState, event screen.Event) screen.Result {
	focus, ok := n.focusItem()

	definition := focus.Node.Screen.Keys()
	required := ok && definition.IsRequired(event.Key)

	if required {
		result := n.focusTick(uiState, event, focus)
		if event.Key.Code != key.ActionEsc {
			return result
		}
	}

	return n.localTick(uiState, event)
}

func (n *Form) localTick(uiState *state.UIState, event screen.Event) screen.Result {
	ky := event.Key

	switch ky.Code {
	case key.ActionEsc:
		n.focused = false
	case key.ActionArrowUp:
		n.cursor = n.decCursor()
	case key.ActionArrowDown:
		n.cursor = n.incCursor()
	case key.ActionArrowLeft:
		n.setCursor(0)
		n.cursor = n.incCursor(0)
	case key.ActionArrowRight:
		items := len(n.items)
		n.setCursor(uint16(items))
		n.cursor = n.decCursor(0)
	case key.ActionEnter:
		n.focused = true
	case key.CustomActionPointer:
		n.pointer = form.NextPointer(n.pointer)
	}

	return screen.ResultFromUIState(uiState)
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

func (n *Form) focusTick(uiState *state.UIState, event screen.Event, focus entry.Entry) screen.Result {
	result := focus.Node.Screen.Tick(uiState, event)

	if result.Node == nil {
		return result

	}

	newItems := make([]entry.Entry, len(n.items))
	copy(newItems, n.items)

	newSteps := make([]pipeline.Transformer, len(n.steps))
	copy(newSteps, n.steps)

	newWrapper := New()
	newWrapper.reference = n.reference
	newWrapper.pointer = n.pointer
	newWrapper.focused = n.focused
	newWrapper.cursor = n.cursor
	newWrapper.steps = newSteps
	newWrapper.items = newItems

	newNode := newWrapper.ToNode()
	result.Node = &newNode

	return result
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

		vm.Kernel.PushLayer(unit, e.Opts...)

		if cvm.Behavior.NeedsPulse {
			vm.Behavior.NeedsPulse = true
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

	return n.applySteps(*vm)
}

func (n *Form) applySteps(vm viewmodel.ViewModel) viewmodel.ViewModel {
	for _, s := range n.steps {
		vm = s(vm)
	}
	return vm
}

func (n *Form) focusItem() (entry.Entry, bool) {
	if n.cursor >= uint16(len(n.items)) {
		return entry.Entry{}, false
	}

	return n.items[n.cursor], true
}
