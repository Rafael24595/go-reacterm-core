package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap/rw"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type CommandRead uint8

const (
	CmdReadNone CommandRead = iota

	CmdReadWriteMode
)

type CommandWrite uint8

const (
	CmdWriteNone CommandWrite = iota

	CmdWriteReadMode

	CmdWriteExecuteAction

	CmdWriteMoveUp
	CmdWriteMoveDown
	CmdWriteMoveLeft
	CmdWriteMoveRight
)

var defaultBindings = rw.Bindings[CommandRead, CommandWrite]{
	Read:  defaultReadBindings,
	Write: defaultWriteBindings,
}

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Write mode", "RET"))

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
	CmdWriteExecuteAction,
	CmdWriteMoveLeft,
	CmdWriteMoveRight,
	CmdWriteMoveUp,
	CmdWriteMoveDown,
}

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Read mode", "ESC")).
	Bind(key.ActionEnter, CmdWriteExecuteAction, key.NewDescriptor("Active selected", "RET")).
	Bind(key.ActionArrowUp, CmdWriteMoveUp).
	Bind(key.ActionArrowDown, CmdWriteMoveDown).
	Bind(key.ActionArrowLeft, CmdWriteMoveLeft).
	Bind(key.ActionArrowRight, CmdWriteMoveRight)
