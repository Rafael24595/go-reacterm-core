package history

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdNone Command = iota

	CmdPrev
	CmdNext
)

var Commands = []Command{
	CmdPrev,
	CmdNext,
}

var defaultBindings = keymap.NewBindings[Command]().
	Bind(key.CustomActionPrev, CmdPrev).
	Bind(key.CustomActionNext, CmdNext)
