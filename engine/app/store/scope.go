package store

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/argument"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Scope struct {
	mu        sync.RWMutex
	timestamp int64
	context   map[string]Argument
}

func newScope(clock clock.Clock) *Scope {
	return &Scope{
		timestamp: clock(),
		context:   make(map[string]Argument),
	}
}

func (n *Scope) Find(key string) (*argument.Argument, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	return &arg.argument, true
}

func (n *Scope) Push(key string, arg Argument) *Scope {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.context[key] = arg

	return n
}

func (n *Scope) Remove(key string) (*argument.Argument, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	arg, ok := n.context[key]
	if !ok {
		return nil, false
	}

	delete(n.context, key)

	return &arg.argument, true
}
