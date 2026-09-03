package store

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Scope struct {
	timestamp int64
	context   map[string]Entry
}

func newScope(clock clock.Clock) *Scope {
	return &Scope{
		timestamp: clock(),
		context:   make(map[string]Entry),
	}
}

func (n *Scope) Find(key string) (*dynamic.Value, bool) {
	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	return &arg.argument, true
}

func (n *Scope) Push(key string, arg Entry) *Scope {
	n.context[key] = arg

	return n
}

func (n *Scope) Remove(key string) (*dynamic.Value, bool) {
	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	delete(n.context, key)

	return &arg.argument, true
}
