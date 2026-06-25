package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
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

	CmdWriteMoveHome
	CmdWriteMoveEnd

	CmdWriteMoveBackward
	CmdWriteMoveForward
	CmdWriteMoveUp
	CmdWriteMoveDown

	CmdWriteDeleteCharBackward
	CmdWriteDeleteWordBackward
	CmdWriteDeleteCharForward
	CmdWriteDeleteWordForward

	CmdWriteUndo
	CmdWriteRedo

	CmdWriteCut
	CmdWriteCopy
	CmdWritePaste

	sysWriteNewLine
	sysWriteRune
)

var defaultBindings = rw.Bindings[CommandRead, CommandWrite]{
	Read:  defaultReadBindings,
	Write: defaultWriteBindings,
}

var CommandsRead = []CommandRead{
	CmdReadWriteMode,
}

var defaultReadBindings = keymap.NewBindings[CommandRead]().
	Bind(key.ActionEnter, CmdReadWriteMode, key.NewDescriptor("Edit mode", "RET"))

var systemWrite = []CommandWrite{
	sysWriteNewLine,
	sysWriteRune,
}

var CommandsWrite = []CommandWrite{
	CmdWriteReadMode,
	CmdWriteMoveHome,
	CmdWriteMoveEnd,
	CmdWriteMoveBackward,
	CmdWriteMoveForward,
	CmdWriteMoveUp,
	CmdWriteMoveDown,
	CmdWriteDeleteCharBackward,
	CmdWriteDeleteWordBackward,
	CmdWriteDeleteCharForward,
	CmdWriteDeleteWordForward,
	CmdWriteUndo,
	CmdWriteRedo,
	CmdWriteCut,
	CmdWriteCopy,
	CmdWritePaste,
}

var systemWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEnter, sysWriteNewLine, key.NewDescriptor("New line", "RET")).
	Bind(key.ActionRune, sysWriteRune)

var defaultWriteBindings = keymap.NewBindings[CommandWrite]().
	Bind(key.ActionEsc, CmdWriteReadMode, key.NewDescriptor("Read mode", "ESC")).
	Bind(key.ActionHome, CmdWriteMoveHome).
	Bind(key.ActionEnd, CmdWriteMoveEnd).
	Bind(key.ActionArrowLeft, CmdWriteMoveBackward).
	Bind(key.ActionArrowRight, CmdWriteMoveForward).
	Bind(key.ActionArrowUp, CmdWriteMoveUp).
	Bind(key.ActionArrowDown, CmdWriteMoveDown).
	Bind(key.ActionBackspace, CmdWriteDeleteCharBackward).
	Bind(key.ActionDeleteBackward, CmdWriteDeleteWordBackward).
	Bind(key.ActionDelete, CmdWriteDeleteCharForward).
	Bind(key.ActionDeleteForward, CmdWriteDeleteWordForward).
	Bind(key.CustomActionUndo, CmdWriteUndo).
	Bind(key.CustomActionRedo, CmdWriteRedo).
	Bind(key.CustomActionCut, CmdWriteCut).
	Bind(key.CustomActionCopy, CmdWriteCopy).
	Bind(key.CustomActionPaste, CmdWritePaste)

var predicates = map[bool]predicate.Predicate{
	false: predicate.Page(),
	true:  predicate.Focus(),
}
