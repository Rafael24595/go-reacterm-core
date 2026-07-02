package clip

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdWriteNone Command = iota

	CmdWriteDec
	CmdWriteInc
)

var CommandsWrite = []Command{
	CmdWriteDec,
	CmdWriteInc,
}

var defaultReadBindings = keymap.NewBindings[Command]()

var defaultWriteBindings = keymap.NewBindings[Command]().
	Bind(key.ActionMinus, CmdWriteDec).
	Bind(key.ActionPlus, CmdWriteInc)
