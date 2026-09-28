package spec

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
)

type argMap = map[ArgKey]dynamic.Value

type args struct {
	items argMap
}

func (a *args) lazyInit() *args {
	if a.items == nil {
		a.items = make(argMap)
	}
	return a
}

// Get retrieves the dynamic Value associated with key, returning zero Value if absent.
func (a *args) Get(key ArgKey) dynamic.Value {
	if a.items == nil {
		var zero dynamic.Value
		return zero
	}

	return a.items[key]
}

// TryGet attempts to retrieve the Value for key, returning the value and existence boolean.
func (a *args) TryGet(key ArgKey) (dynamic.Value, bool) {
	if a.items == nil {
		var zero dynamic.Value
		return zero, false
	}

	v, ok := a.items[key]
	return v, ok
}

// Set stores the specified key-value pair, lazily initializing storage if needed.
func (a *args) Set(key ArgKey, value dynamic.Value) {
	a.lazyInit()
	a.items[key] = value
}

// Delete removes the value for key, returning the deleted value and a boolean indicating if it existed.
func (a *args) Delete(key ArgKey) (dynamic.Value, bool) {
	if a.items == nil {
		var zero dynamic.Value
		return zero, false
	}

	old, ok := a.items[key]
	if !ok {
		var zero dynamic.Value
		return zero, false
	}

	delete(a.items, key)
	return old, true
}

// Copy merges all key-value pairs from src into current storage, returning the underlying map.
func (a *args) Copy(src args) argMap {
	a.lazyInit()

	maps.Copy(a.items, src.Items())
	return a.items
}

// Clone creates a deep copy of the args instance, duplicating internal map entries.
func (a *args) Clone() args {
	a.lazyInit()

	return args{
		items: maps.Clone(a.items),
	}
}

// Items returns the internal argument map representation, initializing it lazily if unallocated.
func (a *args) Items() argMap {
	a.lazyInit()
	return a.items
}
