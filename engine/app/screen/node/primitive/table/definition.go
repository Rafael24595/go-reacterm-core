package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
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

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Read mode", "RET"))

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
	CmdWriteExecuteAction,
	CmdWriteMoveLeft,
	CmdWriteMoveRight,
	CmdWriteMoveUp,
	CmdWriteMoveDown,
}

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Write mode", "ESC")).
	Bind(key.ActionEnter, CmdWriteExecuteAction, key.NewDescriptor("Active selected", "RET")).
	Bind(key.ActionArrowUp, CmdWriteMoveUp).
	Bind(key.ActionArrowDown, CmdWriteMoveDown).
	Bind(key.ActionArrowLeft, CmdWriteMoveLeft).
	Bind(key.ActionArrowRight, CmdWriteMoveRight)

type bindings struct {
	read  *keymap.Bindings[CommandRead]
	write *keymap.Bindings[CommandWrite]
}

var defaultBindings = bindings{
	read:  defaultReadBindings,
	write: defaultWriteBindings,
}

type definition struct {
	read  screen.Definition
	write screen.Definition
}

func emptyDefinition() definition {
	return definition{
		read:  screen.EmptyDefinition(),
		write: screen.EmptyDefinition(),
	}
}

func definitionFromBindings(bindings bindings) definition {
	return definition{
		read:  keymap.BindingsToDefinition(bindings.read),
		write: keymap.BindingsToDefinition(bindings.write),
	}
}

func (d definition) get(write bool) screen.Definition {
	if write {
		return d.write
	}
	return d.read
}
