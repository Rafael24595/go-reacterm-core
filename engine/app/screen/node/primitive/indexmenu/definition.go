package indexmenu

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

	CmdSwitchPointer
)

var Commands = []Command{
	CmdExecuteAction,
	CmdPrevOption,
	CommandNextOption,
	CmdFirstOption,
	CmdLastOption,
	CmdSwitchPointer,
}

var defaultBindings = keymap.NewBindings[Command]().
	Bind(key.ActionEnter, CmdExecuteAction, key.NewDescriptor("Accept", "RET")).
	Bind(key.ActionArrowLeft, CmdFirstOption,  key.NewDescriptor("←", "Move first")).
	Bind(key.ActionArrowRight, CmdLastOption, key.NewDescriptor("→", "Move last")).
	Bind(key.ActionArrowUp, CmdPrevOption).
	Bind(key.ActionArrowDown, CommandNextOption).
	Bind(key.CustomActionPointer, CmdSwitchPointer)
