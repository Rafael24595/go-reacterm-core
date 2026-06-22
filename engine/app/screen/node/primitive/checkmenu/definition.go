package checkmenu

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

	CmdWriteSwitchState

	CmdWritePrevOption
	CmdWriteNextOption

	CmdWriteFirstOption
	CmdWriteLastOption
)

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Read mode", "RET"))

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
	CmdWriteSwitchState,
	CmdWritePrevOption,
	CmdWriteNextOption,
	CmdWriteFirstOption,
	CmdWriteLastOption,
}

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Write mode", "ESC")).
	Bind(key.ActionEnter, CmdWriteSwitchState, key.NewDescriptor("Switch state", "RET")).
	Bind(key.ActionArrowUp, CmdWriteFirstOption, key.NewDescriptor("↑", "Move first")).
	Bind(key.ActionArrowDown, CmdWriteLastOption, key.NewDescriptor("↓", "Move last")).
	Bind(key.ActionArrowLeft, CmdWritePrevOption).
	Bind(key.ActionArrowRight, CmdWriteNextOption)

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
