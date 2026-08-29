package rw

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

// Definition holds separate screen keymap definitions for read-only and write/editing states.
type Definition struct {
	// Read represents the keymap definition applied when the component is in read mode.
	Read  screen.Definition
	// Write represents the keymap definition applied when the component is in write mode.
	Write screen.Definition
}

// EmptyDefinition creates and returns a Definition initialized with empty screen definitions for both Read and Write states.
func EmptyDefinition() Definition {
	return Definition{
		Read:  screen.EmptyDefinition(),
		Write: screen.EmptyDefinition(),
	}
}

// Get returns either the Write or Read screen.Definition based on the provided write flag.
func (d Definition) Get(write bool) screen.Definition {
	if write {
		return d.Write
	}
	return d.Read
}
