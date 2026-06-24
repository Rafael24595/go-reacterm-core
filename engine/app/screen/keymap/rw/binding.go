package rw

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"

type Bindings[T, K keymap.Command] struct {
	Read  *keymap.Bindings[T]
	Write *keymap.Bindings[K]
}

func DefinitionFromBindings[T, K keymap.Command](bindings Bindings[T, K]) Definition {
	return Definition{
		Read:  keymap.BindingsToDefinition(bindings.Read),
		Write: keymap.BindingsToDefinition(bindings.Write),
	}
}
