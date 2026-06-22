package modalmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdNone Command = iota

	CmdExecuteAction

	CmdPrevOption
	CommandNextOption

	CmdFirstOption
	CmdLastOption
)

var Commands = []Command{
	CmdExecuteAction,
	CmdPrevOption,
	CommandNextOption,
	CmdFirstOption,
	CmdLastOption,
}

var defaultBindings = keymap.NewBindings[Command]().
	Bind(key.ActionEnter, CmdExecuteAction, key.NewDescriptor("Active selected", "RET")).
	Bind(key.ActionArrowUp, CmdFirstOption, key.NewDescriptor("↑", "Move first")).
	Bind(key.ActionArrowDown, CmdLastOption, key.NewDescriptor("↓", "Move last")).
	Bind(key.ActionArrowLeft, CmdPrevOption).
	Bind(key.ActionArrowRight, CommandNextOption)
