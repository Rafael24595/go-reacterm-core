package rw

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"

// Bindings groups distinct keymap bindings for read-only and write/editing states.
type Bindings[T, K keymap.Command] struct {
	// Read holds the key bindings active during read-only interactions.
	Read  *keymap.Bindings[T]
	// Write holds the key bindings active during write/editing interactions.
	Write *keymap.Bindings[K]
}

// DefinitionFromBindings converts a dual-mode Bindings struct into a Definition containing screen definitions for both Read and Write states.
func DefinitionFromBindings[T, K keymap.Command](bindings Bindings[T, K]) Definition {
	return Definition{
		Read:  keymap.BindingsToDefinition(bindings.Read),
		Write: keymap.BindingsToDefinition(bindings.Write),
	}
}
