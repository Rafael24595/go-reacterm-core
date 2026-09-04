package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

// Scope manages a isolated set of Arguments identified by string keys.
// State synchronization is handled at the root Store level.
type Scope struct {
	clock     clock.Clock
	timestamp int64
	context   map[string]Entry
}

func newScope(clock clock.Clock) *Scope {
	return &Scope{
		clock:     clock,
		timestamp: clock(),
		context:   make(map[string]Entry),
	}
}

// Find retrieves an argument's dynamic value by key within the scope.
// Returns the pointer to the dynamic value and true if present, or nil and false otherwise.
func (n *Scope) Find(key string) (*dynamic.Value, bool) {
	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	return &arg.value, true
}

// Push inserts or overwrites an Argument in the scope and updates the scope's timestamp.
func (n *Scope) Push(key string, arg Entry) *Scope {
	n.timestamp = n.clock()
	n.context[key] = arg

	return n
}

// Remove deletes an Argument from the scope by key and updates the scope's timestamp.
// Returns the removed dynamic value pointer and true if found, or nil and false otherwise.
func (n *Scope) Remove(key string) (*dynamic.Value, bool) {
	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	n.timestamp = n.clock()
	delete(n.context, key)

	return &arg.value, true
}
