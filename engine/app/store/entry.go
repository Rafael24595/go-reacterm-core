package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

// Entry represents an immutable record within a scope, wrapping its dynamic value
// alongside creation/modification metadata.
type Entry struct {
	timestamp int64
	value     dynamic.Value
}

func newEntry(clk clock.Clock, arg any) Entry {
	return Entry{
		timestamp: clk(),
		value:     dynamic.From(arg),
	}
}

// Timestamp returns the Unix timestamp when this entry was created or updated.
func (e Entry) Timestamp() int64 {
	return e.timestamp
}

// Value returns the underlying dynamic value wrapper.
func (e Entry) Value() dynamic.Value {
	return e.value
}
