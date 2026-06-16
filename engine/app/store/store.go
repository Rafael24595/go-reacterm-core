package store

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Updater[T any] func(*T)

type Store struct {
	mu     sync.RWMutex
	clock  clock.Clock
	scopes map[string]*Scope
}

func New() *Store {
	return &Store{
		clock:  clock.UnixMilliClock,
		scopes: make(map[string]*Scope),
	}
}

func (n *Store) Find(scope string, key string) (*commons.Argument, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()

	ctx, ok := n.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Find(key)
}

func (n *Store) Push(scope string, key string, arg any) *Store {
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, ok := n.scopes[scope]
	if !ok {
		ctx = newScope(n.clock)
	}

	n.scopes[scope] = ctx.Push(key,
		newArgument(n.clock, arg),
	)

	return n
}

func (n *Store) RemoveScope(scope string) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, ok := n.scopes[scope]
	if !ok {
		return false
	}

	delete(n.scopes, scope)

	return true
}

func (n *Store) RemoveArgument(scope, key string) (*commons.Argument, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()

	ctx, ok := n.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Remove(key)
}

func (n *Store) RetainOnly(scopes set.Set[string]) *Store {
	n.mu.Lock()
	items := make([]string, 0)

	for scope := range n.scopes {
		if !scopes.Has(scope) {
			items = append(items, scope)
		}
	}

	n.mu.Unlock()

	for _, name := range items {
		n.RemoveScope(name)
	}

	return n
}

func Find[T any](
	c *Store,
	scope string,
	key Key[T],
) (T, bool) {
	arg, ok := c.Find(scope, key.Code())
	if ok {
		return commons.Map[T](*arg)
	}

	var zero T
	return zero, false
}

func Push[T any](
	c *Store,
	scope string,
	key Key[T],
	arg T,
) *Store {
	return c.Push(scope, key.Code(), arg)
}

func Update[T any](
	c *Store,
	scope string,
	key Key[T],
	updater Updater[T],
) (T, bool) {
	arg, ok := c.Find(scope, key.Code())
	if !ok {
		var zero T
		return zero, false
	}

	value, ok := commons.Map[T](*arg)
	if !ok {
		var zero T
		return zero, false
	}

	updater(&value)
	Push(c, scope, key, value)

	return value, true
}

func Remove[T any](
	c *Store,
	scope string,
	key Key[T],
) (T, bool) {
	arg, ok := c.RemoveArgument(scope, key.Code())
	if ok {
		return commons.Map[T](*arg)
	}

	var zero T
	return zero, false
}
