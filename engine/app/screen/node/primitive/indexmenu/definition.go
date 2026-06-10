package indexmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var definition = screen.DefinitionFromActions(
	[]key.Action{
		key.ActionEnter,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionPointer,
	}...,
)
