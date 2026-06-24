package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap/rw"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type CommandWrite uint8

const (
	CmdWriteNone CommandWrite = iota

	CmdWriteReadMode
)

type CommandRead uint8

const (
	CmdReadNone CommandRead = iota

	CmdReadWriteMode

	CmdReadSwitchPointer

	CmdReadPrevOption
	CmdReadNextOption

	CmdReadFirstOption
	CmdReadLastOption
)

var defaultBindings = rw.Bindings[CommandRead, CommandWrite]{
	Write: defaultWriteBindings,
	Read:  defaultReadBindings,
}

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
}

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Write mode", "RET"))

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
	CmdReadPrevOption,
	CmdReadNextOption,
	CmdReadFirstOption,
	CmdReadLastOption,
	CmdReadSwitchPointer,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Read mode", "ESC")).
	Bind(key.ActionArrowLeft, CmdReadFirstOption, key.NewDescriptor("←", "Move first")).
	Bind(key.ActionArrowRight, CmdReadLastOption, key.NewDescriptor("→", "Move last")).
	Bind(key.ActionArrowUp, CmdReadPrevOption).
	Bind(key.ActionArrowDown, CmdReadNextOption).
	Bind(key.CustomActionPointer, CmdReadSwitchPointer)
