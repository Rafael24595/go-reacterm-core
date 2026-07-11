package store

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
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

func (s *Store) Find(scope string, key string) (*dynamic.Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Find(key)
}

func (s *Store) Push(scope string, key string, arg any) *Store {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		ctx = newScope(s.clock)
	}

	s.scopes[scope] = ctx.Push(key,
		newArgument(s.clock, arg),
	)

	return s
}

func (s *Store) RemoveScope(scope string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.scopes[scope]
	if !ok {
		return false
	}

	delete(s.scopes, scope)

	return true
}

func (s *Store) RemoveArgument(scope, key string) (*dynamic.Value, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Remove(key)
}

func (s *Store) RetainOnly(scopes set.Set[string]) *Store {
	s.mu.Lock()
	items := make([]string, 0)

	for scope := range s.scopes {
		if !scopes.Has(scope) {
			items = append(items, scope)
		}
	}

	s.mu.Unlock()

	for _, name := range items {
		s.RemoveScope(name)
	}

	return s
}

func Find[T any](
	store *Store,
	scope string,
	key Key[T],
) (T, bool) {
	arg, ok := store.Find(scope, key.Code())
	if ok {
		return dynamic.Map[T](*arg)
	}

	var zero T
	return zero, false
}

func Push[T any](
	store *Store,
	scope string,
	key Key[T],
	arg T,
) *Store {
	return store.Push(scope, key.Code(), arg)
}

func Update[T any](
	store *Store,
	scope string,
	key Key[T],
	updater Updater[T],
) (T, bool) {
	arg, ok := store.Find(scope, key.Code())
	if !ok {
		var zero T
		return zero, false
	}

	value, ok := dynamic.Map[T](*arg)
	if !ok {
		var zero T
		return zero, false
	}

	updater(&value)
	Push(store, scope, key, value)

	return value, true
}

func Upsert[T any](
	store *Store,
	scope string,
	key Key[T],
	updater Updater[T],
) (T, bool) {
	value, _ := Find(store, scope, key)

	updater(&value)
	Push(store, scope, key, value)

	return value, true
}

func Remove[T any](
	store *Store,
	scope string,
	key Key[T],
) (T, bool) {
	arg, ok := store.RemoveArgument(scope, key.Code())
	if ok {
		return dynamic.Map[T](*arg)
	}

	var zero T
	return zero, false
}
