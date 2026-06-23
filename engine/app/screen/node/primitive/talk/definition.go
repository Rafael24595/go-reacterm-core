package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
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

	CmdWritePrevOption
	CmdWriteNextOption

	CmdWriteFirstOption
	CmdWriteLastOption

	CmdWriteSwitchPointer
)

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Write mode", "RET"))

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
	CmdWritePrevOption,
	CmdWriteNextOption,
	CmdWriteFirstOption,
	CmdWriteLastOption,
	CmdWriteSwitchPointer,
}

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Read mode", "ESC")).
	Bind(key.ActionArrowLeft, CmdWriteFirstOption, key.NewDescriptor("←", "Move first")).
	Bind(key.ActionArrowRight, CmdWriteLastOption, key.NewDescriptor("→", "Move last")).
	Bind(key.ActionArrowUp, CmdWritePrevOption).
	Bind(key.ActionArrowDown, CmdWriteNextOption).
	Bind(key.CustomActionPointer, CmdWriteSwitchPointer)

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

var predicates = map[bool]predicate.Predicate{
	false: predicate.Page(),
	true:  predicate.Focus(),
}

var actions = map[bool]action.Action{
	false: action.Scroll(),
	true:  action.Paged(),
}
