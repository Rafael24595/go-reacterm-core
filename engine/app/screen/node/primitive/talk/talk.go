package talk

import (
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/talk"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

const Name = "talk"

type Talk struct {
	reference  string
	navigation bool
	pointer    uint8
	owner      string
	messages   []chat.Message
	cursor     uint16
}

func New() *Talk {
	return &Talk{
		reference: Name,
		owner:     "",
		messages:  make([]chat.Message, 0),
		cursor:    0,
	}
}

func (n *Talk) SetName(name string) *Talk {
	n.reference = name
	return n
}

func (n *Talk) SetOwner(owner string) *Talk {
	n.owner = owner
	return n
}

func (n *Talk) AddMessage(message ...chat.Message) *Talk {
	n.messages = append(n.messages, message...)
	return n
}

func (n *Talk) SetCursor(cursor uint16) *Talk {
	maxIdx := math.SubClampZeroAs[int, uint16](len(n.messages), 1)
	n.cursor = math.Clamp(cursor, 0, maxIdx)
	return n
}

func (n *Talk) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Init(n.init).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *Talk) init(uiState state.UIState) {
	n.loadFromStack(uiState)
}

func (n *Talk) loadFromStack(uiState state.UIState) {
	if cursor, ok := KeyCursor.Get(
		uiState.Store,
		n.reference,
	); ok {
		n.cursor = cursor
	}

	if messages, ok := KeyMessages.Get(
		uiState.Store,
		n.reference,
	); ok {
		n.messages = messages
	}
}

func (n *Talk) keys() screen.Definition {
	return definitions[n.navigation]
}

func (n *Talk) tick(uiState *state.UIState, event screen.Event) screen.Result {
	if n.navigation {
		return n.tickNavigation(uiState, event)
	}

	switch event.Key.Code {
	case key.ActionEnter:
		n.navigation = true
	}

	return screen.ResultFromUIState(uiState)
}

func (n *Talk) tickNavigation(uiState *state.UIState, event screen.Event) screen.Result {
	size := uint16(len(n.messages))
	if size == 0 {
		return screen.ResultFromUIState(uiState)
	}

	switch event.Key.Code {
	case key.ActionEsc:
		n.navigation = false
	case key.ActionArrowUp:
		n.cursor = (n.cursor + size - 1) % size
		n.tickToStack(uiState)
	case key.ActionArrowDown:
		n.cursor = (n.cursor + 1) % size
		n.tickToStack(uiState)
	case key.ActionArrowLeft:
		n.cursor = 0
		n.tickToStack(uiState)
	case key.ActionArrowRight:
		optsLen := uint16(len(n.messages))
		optsLen = math.SubClampZero(optsLen, 1)
		n.cursor = min(optsLen, n.cursor+1)
		n.tickToStack(uiState)
	case key.CustomActionPointer:
		n.pointer = talk.NextPointer(n.pointer)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *Talk) tickToStack(uiState *state.UIState) {
	KeyCursor.Set(
		uiState.Store,
		n.reference,
		n.cursor,
	)

	KeyMessages.Set(
		uiState.Store,
		n.reference,
		n.messages,
	)
}

func (n *Talk) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStack(uiState)

	pointer := n.pointerProvider()

	indexmenu := talk.New().
		Navigation(n.navigation).
		Pointer(pointer).
		SetOwner(n.owner).
		AddMessage(n.messages...).
		SetCursor(n.cursor)

	vm.Kernel.Push(
		indexmenu.ToUnit(),
	)

	if n.navigation {
		index := math.Clamp(
			n.cursor, 0, uint16(len(n.messages)),
		)

		text := fmt.Sprintf(
			"%d - %s", n.cursor, n.messages[index].Owner,
		)

		vm.Footer.Push(
			inputline.FromString(text),
		)
	}

	vm.Pager.Action = actions[n.navigation]

	vm.Pager.SetPredicate(
		predicates[n.navigation],
	)

	return *vm
}

func (n *Talk) pointerProvider() talk.PointerProvider {
	if n.navigation {
		return talk.FindPointer(n.pointer)
	}
	return talk.NoneProvider
}
