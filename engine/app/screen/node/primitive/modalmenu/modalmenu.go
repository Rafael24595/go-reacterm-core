package modalmenu

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/modal"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "modal_menu"

type ModalMenu struct {
	reference  string
	loaded     bool
	bindings   *keymap.Bindings[Command]
	definition screen.Definition
	text       []text.Line
	options    []input.MenuOption
	cursor     uint16
}

func New() *ModalMenu {
	return &ModalMenu{
		reference:  Name,
		loaded:     false,
		bindings:   defaultBindings,
		definition: screen.EmptyDefinition(),
		text:       make([]text.Line, 0),
		options:    make([]input.MenuOption, 0),
		cursor:     0,
	}
}

func (n *ModalMenu) SetName(name string) *ModalMenu {
	n.reference = name
	return n
}

func (n *ModalMenu) WithBindings(overrides *keymap.Bindings[Command]) *ModalMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings = n.bindings.Overlay(overrides)
	return n
}

func (n *ModalMenu) AddText(text ...text.Line) *ModalMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.text = append(n.text, text...)
	return n
}

func (n *ModalMenu) AddOptions(options ...input.MenuOption) *ModalMenu {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	n.options = append(n.options, options...)
	return n
}

func (n *ModalMenu) SetCursor(cursor uint16) *ModalMenu {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.options), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *ModalMenu) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *ModalMenu) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loaded = true

	n.loadFromStore(uiState)
	n.definition = keymap.BindingsToDefinition(n.bindings)
}

func (n *ModalMenu) loadFromStore(uiState state.UIState) {
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

func (n *ModalMenu) keys() screen.Definition {
	return n.definition
}

func (n *ModalMenu) tick(uiState *state.UIState, event screen.Event) screen.Result {

	switch n.bindings.Command(event.Key.Code) {
	case CmdExecuteAction:
		n.tickToStore(uiState)
		return n.actionEnter()
	case CmdPrevOption:
		n.cursor = math.SubClampZero(n.cursor, 1)
		n.tickToStore(uiState)
	case CommandNextOption:
		last := math.SubClampZeroAs[int, uint16](len(n.options), 1)
		n.cursor = min(last, n.cursor+1)
		n.tickToStore(uiState)
	case CmdFirstOption:
		n.cursor = 0
		n.tickToStore(uiState)
	case CmdLastOption:
		n.cursor = math.SubClampZeroAs[int, uint16](len(n.options), 1)
		n.tickToStore(uiState)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *ModalMenu) tickToStore(uiState *state.UIState) {
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

func (n *ModalMenu) actionEnter() screen.Result {
	node := n.options[n.cursor].Action()
	return screen.ResultFromNode(&node)
}

func (n *ModalMenu) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStore(uiState)

	frags := input.FragmentFromMenuOption(n.options...)

	modal := modal.New().
		AddText(n.text...).
		AddOptions(frags...).
		SetCursor(n.cursor).
		ToUnit()

	vm.Kernel.Push(modal)

	return *vm
}
