package key

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
)

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
	ActionArrowUp:    NewDescriptor("Move up", "↑"),
	ActionArrowDown:  NewDescriptor("Move down", "↓"),
	ActionArrowLeft:  NewDescriptor("Move left", "←"),
	ActionArrowRight: NewDescriptor("Move right", "→"),
	ActionHome:       NewDescriptor("Line start", "HOME", "^A"),
	ActionEnd:        NewDescriptor("Line end", "END", "^E"),

	ActionEnter: NewDescriptor("New line/Accept", "RET"),
	ActionTab:   NewDescriptor("Next field", "TAB"),
	ActionEsc:   NewDescriptor("Back/Cancel", "ESC"),
	ActionExit:  NewDescriptor("Exit", "^C"),

	ActionBackspace:      NewDescriptor("Delete char", "BS"),
	ActionDelete:         NewDescriptor("Delete forward", "DEL"),
	ActionDeleteBackward: NewDescriptor("Delete word", "^W"),
	ActionDeleteForward:  NewDescriptor("Delete word fwd", "^D"),

	ActionPageUp:   NewDescriptor("⇞", "Prev page"),
	ActionPageDown: NewDescriptor("⇟", "Next page"),

	CustomActionUndo:  NewDescriptor("Undo", "M-z"),
	CustomActionRedo:  NewDescriptor("Redo", "M-y"),
	CustomActionHelp:  NewDescriptor("Help", "M-h"),
	CustomActionBack:  NewDescriptor("Back", "M-b"),
	CustomActionCut:   NewDescriptor("Cut", "M-x"),
	CustomActionCopy:  NewDescriptor("Copy", "M-c"),
	CustomActionPaste: NewDescriptor("Paste", "M-v"),

	CustomActionPointer: NewDescriptor("Switch pointer", "M-s"),

	ActionRune: NewDescriptor("Text", "Text"),
}

func ResolveDescriptors(actions ...Action) *dict.LinkedMap[Action, Descriptor] {
	return ResolveDescriptorsWithDefaults(nil, actions...)
}

func ResolveDescriptorsWithDefaults(
	defaults map[Action]Descriptor,
	actions ...Action,
) *dict.LinkedMap[Action, Descriptor] {
	help := dict.NewLinkedMap[Action, Descriptor]()
	for _, a := range actions {
		if action := resolveDescriptor(defaults, a); action != nil {
			help.Set(a, *action)
		}
	}

	return help
}

func FindDescriptor(action Action) *Descriptor {
	return resolveDescriptor(make(map[Action]Descriptor), action)
}

func resolveDescriptor(
	defaults map[Action]Descriptor,
	action Action,
) *Descriptor {
	if action == ActionAll {
		return nil
	}

	if defaults != nil {
		if field, exists := defaults[action]; exists {
			return &field
		}
	}

	if str, exist := descriptors[action]; exist {
		return &str
	}

	assert.Unreachable("unhandled action: %d", action)

	descriptor := NewDescriptor("Unknown action", "???")
	return &descriptor
}
