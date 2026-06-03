package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var read_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit text"},
	},
	[]key.Action{
		key.ActionEnter,
	},
)

var write_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Save & Quit"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "New line"},
	},
	[]key.Action{
		key.ActionEsc,
		key.ActionHome,
		key.ActionEnd,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionBackspace,
		key.ActionDeleteBackward,
		key.ActionDelete,
		key.ActionDeleteForward,
		key.ActionEnter,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionUndo,
		key.CustomActionRedo,
		key.CustomActionCut,
		key.CustomActionCopy,
		key.CustomActionPaste,
		key.ActionRune,
	},
)
