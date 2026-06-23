package history

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdNone Command = iota

	CmdBack
)

var Commands = []Command{
	CmdBack,
}

var defaultBindings = keymap.NewBindings[Command]().
	Bind(key.CustomActionBack, CmdBack)
