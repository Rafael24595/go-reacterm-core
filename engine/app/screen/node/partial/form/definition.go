package form

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
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

type bindings struct {
	write *keymap.Bindings[CommandWrite]
	read  *keymap.Bindings[CommandRead]
}

var defaultBindings = bindings{
	write: defaultWriteBindings,
	read:  defaultReadBindings,
}

type definition struct {
	write screen.Definition
	read  screen.Definition
}

func emptyDefinition() definition {
	return definition{
		write: screen.EmptyDefinition(),
		read:  screen.EmptyDefinition(),
	}
}

func definitionFromBindings(bindings bindings) definition {
	return definition{
		write: keymap.BindingsToDefinition(bindings.write),
		read:  keymap.BindingsToDefinition(bindings.read),
	}
}

func (d definition) get(write bool) screen.Definition {
	if write {
		return d.write
	}
	return d.read
}
