package help

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdNone Command = iota

	CmdSwitchDisplay
)

var Commands = []Command{
	CmdSwitchDisplay,
}

var defaultBindings = keymap.NewBindings[Command]().
	Bind(key.CustomActionHelp, CmdSwitchDisplay)
