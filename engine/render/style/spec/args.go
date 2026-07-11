package spec

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
)

type argMap = map[ArgKey]dynamic.Value

type args struct {
	items argMap
}

func (a *args) Get(key ArgKey) dynamic.Value {
	if a.items == nil {
		var zero dynamic.Value
		return zero
	}

	return a.items[key]
}

func (a *args) TryGet(key ArgKey) (dynamic.Value, bool) {
	if a.items == nil {
		var zero dynamic.Value
		return zero, false
	}

	v, ok := a.items[key]
	return v, ok
}

func (a *args) Set(key ArgKey, value dynamic.Value) {
	if a.items == nil {
		a.items = make(argMap)
	}
	a.items[key] = value
}

func (a *args) Delete(key ArgKey) (dynamic.Value, bool) {
	if a.items == nil {
		var zero dynamic.Value
		return zero, false
	}

	old := a.items[key]
	delete(a.items, key)

	return old, true
}

func (a *args) Copy(src args) argMap {
	if a.items == nil {
		a.items = make(argMap)
	}
	maps.Copy(a.items, src.Items())
	return a.items
}

func (a *args) Items() argMap {
	if a.items == nil {
		a.items = make(argMap)
	}
	return a.items
}
