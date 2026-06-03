package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var read_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Edit mode"},
	},
	[]key.Action{
		key.ActionEnter,
	},
)

var write_definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEsc:   {Code: []string{"ESC"}, Detail: "Write Mode"},
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.Action{
		key.ActionEsc,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
	},
)
