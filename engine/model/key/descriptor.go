package key

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
)

type DescriptorResolver func(action Action) *Descriptor
type DescriptorsResolver func(action ...Action) *dict.LinkedMap[Action, Descriptor]

type Descriptor struct {
	Code   []string
	Detail string
}

func NewDescriptor(detail string, codes ...string) Descriptor {
	return Descriptor{
		Code:   codes,
		Detail: detail,
	}
}

var descriptors = map[Action]Descriptor{
	ActionRune: NewDescriptor("Text", "Text"),

	ActionEsc:  NewDescriptor("Back/Cancel", "ESC"),
	ActionExit: NewDescriptor("Exit", "^C"),

	ActionDeleteBackward: NewDescriptor("Delete word", "^W"),
	ActionDeleteForward:  NewDescriptor("Delete word fwd", "^D"),

	ActionTab:       NewDescriptor("Next field", "TAB"),
	ActionEnter:     NewDescriptor("New line/Accept", "RET"),
	ActionBackspace: NewDescriptor("Delete char", "BS"),

	ActionArrowUp:    NewDescriptor("Move up", "↑"),
	ActionArrowDown:  NewDescriptor("Move down", "↓"),
	ActionArrowLeft:  NewDescriptor("Move left", "←"),
	ActionArrowRight: NewDescriptor("Move right", "→"),

	ActionHome:   NewDescriptor("Line start", "HOME", "^A"),
	ActionEnd:    NewDescriptor("Line end", "END", "^E"),
	ActionDelete: NewDescriptor("Delete forward", "DEL"),

	ActionPageUp:   NewDescriptor("⇞", "Prev page"),
	ActionPageDown: NewDescriptor("⇟", "Next page"),

	CustomActionHelp: NewDescriptor("Help", "M-h"),

	CustomActionPrev: NewDescriptor("Prev", "M-p"),
	CustomActionNext: NewDescriptor("Prev", "M-n"),

	CustomActionUndo: NewDescriptor("Undo", "M-z"),
	CustomActionRedo: NewDescriptor("Redo", "M-y"),

	CustomActionCut:   NewDescriptor("Cut", "M-x"),
	CustomActionCopy:  NewDescriptor("Copy", "M-c"),
	CustomActionPaste: NewDescriptor("Paste", "M-v"),

	CustomActionPointer: NewDescriptor("Switch pointer", "M-s"),
}

func ResolveDescriptor(action Action) *Descriptor {
	if action == ActionAll {
		return nil
	}

	if str, exist := descriptors[action]; exist {
		return &str
	}

	assert.Unreachable("unhandled action: %d", action)

	descriptor := NewDescriptor("Unknown action", "???")
	return &descriptor
}

func ResolveDescriptors(actions ...Action) *dict.LinkedMap[Action, Descriptor] {
	help := dict.NewLinkedMap[Action, Descriptor]()
	for _, a := range actions {
		if action := ResolveDescriptor(a); action != nil {
			help.Set(a, *action)
		}
	}

	return help
}
