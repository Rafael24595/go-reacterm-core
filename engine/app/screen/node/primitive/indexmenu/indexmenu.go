package indexmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/indexmenu"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

const Name = "index_menu"

type IndexMenu struct {
	reference string
	pointer   uint8
	meta      marker.IndexMeta
	options   []input.MenuOption
	cursor    uint16
}

func New() *IndexMenu {
	return &IndexMenu{
		reference: Name,
		pointer:   0,
		meta:      marker.HyphenIndex,
		options:   make([]input.MenuOption, 0),
		cursor:    0,
	}
}

func (n *IndexMenu) SetName(name string) *IndexMenu {
	n.reference = name
	return n
}

func (n *IndexMenu) SetMeta(meta marker.IndexMeta) *IndexMenu {
	n.meta = meta
	return n
}

func (n *IndexMenu) AddOptions(options ...input.MenuOption) *IndexMenu {
	n.options = append(n.options,
		input.NormalizeMenuOptions(options...)...,
	)
	return n
}

func (n *IndexMenu) SetCursor(cursor uint16) *IndexMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *IndexMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Init(n.init).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *IndexMenu) init(uiState state.UIState) {
	n.loadFromStore(uiState)
}

func (n *IndexMenu) loadFromStore(uiState state.UIState) {
	option, ok := KeyActive.Get(
		uiState.Store,
		n.reference,
	)

	if !ok {
		return
	}

	for i, o := range n.options {
		if o.Id == option {
			n.cursor = uint16(i)
			break
		}
	}
}

func (n *IndexMenu) keys() screen.Definition {
	return definition
}

func (n *IndexMenu) tick(uiState *state.UIState, event screen.Event) screen.Result {
	size := uint16(len(n.options))
	if size == 0 {
		return screen.EmptyResult()
	}

	switch event.Key.Code {
	case key.ActionArrowUp:
		n.cursor = (n.cursor + size - 1) % size
		n.tickToStore(uiState)
	case key.ActionArrowDown:
		n.cursor = (n.cursor + 1) % size
		n.tickToStore(uiState)
	case key.ActionArrowLeft:
		n.cursor = 0
		n.tickToStore(uiState)
	case key.ActionArrowRight:
		optsLen := uint16(len(n.options))
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = min(optsLen, n.cursor+1)
		n.tickToStore(uiState)
	case key.ActionEnter:
		n.tickToStore(uiState)
		return n.actionEnter()
	case key.CustomActionPointer:
		n.pointer = indexmenu.NextPointer(n.pointer)
	}

	return screen.EmptyResult()
}

func (n *IndexMenu) tickToStore(uiState *state.UIState) {
	if n.cursor >= uint16(len(n.options)) {
		KeyActive.Delete(
			uiState.Store,
			n.reference,
		)
		return
	}

	KeyActive.Set(
		uiState.Store,
		n.reference,
		n.options[n.cursor].Id,
	)
}

func (n *IndexMenu) actionEnter() screen.Result {
	node := n.options[n.cursor].Action()
	return screen.ResultFromNode(&node)
}

func (n *IndexMenu) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStore(uiState)

	frags := input.FragmentFromMenuOption(n.options...)

	pointer := indexmenu.FindPointer(n.pointer)

	indexmenu := indexmenu.New(frags).
		Pointer(pointer).
		Meta(n.meta).
		Cursor(n.cursor)

	vm.Kernel.Push(
		indexmenu.ToUnit(),
	)

	index := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	option := min(index, n.cursor)
	text := n.options[option].Label.Text

	vm.Footer.Push(
		inputline.FromString(text),
	)

	vm.Pager.SetPredicate(
		predicate.Focus(),
	)

	return *vm
}
