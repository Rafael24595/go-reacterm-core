package spec

import (
	"maps"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/argument"
)

type argMap = map[ArgKey]argument.Argument

type args struct {
	items argMap
}

func (a *args) Get(key ArgKey) argument.Argument {
	if a.items == nil {
		var zero argument.Argument
		return zero
	}

	return a.items[key]
}

func (a *args) TryGet(key ArgKey) (argument.Argument, bool) {
	if a.items == nil {
		var zero argument.Argument
		return zero, false
	}

	v, ok := a.items[key]
	return v, ok
}

func (a *args) Set(key ArgKey, value argument.Argument) {
	if a.items == nil {
		a.items = make(argMap)
	}
	a.items[key] = value
}

func (a *args) Delete(key ArgKey) (argument.Argument, bool) {
	if a.items == nil {
		var zero argument.Argument
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
