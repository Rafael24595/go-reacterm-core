package modalmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var definition = screen.NewDefinition(
	map[key.Action]key.Descriptor{
		key.ActionEnter: {Code: []string{"RET"}, Detail: "Active selected"},
	},
	[]key.Action{
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionBack,
	},
)
