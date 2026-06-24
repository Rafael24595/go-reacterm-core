package rw

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Definition struct {
	Read  screen.Definition
	Write screen.Definition
}

func EmptyDefinition() Definition {
	return Definition{
		Read:  screen.EmptyDefinition(),
		Write: screen.EmptyDefinition(),
	}
}

func (d Definition) Get(write bool) screen.Definition {
	if write {
		return d.Write
	}
	return d.Read
}
