package store

import (
	"sync"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

// Updater defines a function signature that receives a pointer to a value of type T,
// allowing in-place state modifications.
type Updater[T any] func(*T)

// Store represents a thread-safe, multi-scope key-value store.
// It manages isolated state spaces (scopes) and synchronizes access using a Read-Write Mutex.
type Store struct {
	mu     sync.RWMutex
	clock  clock.Clock
	scopes map[string]*Scope
}

// New creates and initializes a new thread-safe Store instance
// configured with the default Unix millisecond clock.
func New() *Store {
	return &Store{
		clock:  clock.UnixMilliClock,
		scopes: make(map[string]*Scope),
	}
}

// Find retrieves an untyped dynamic value by key within a given scope.
// Returns the dynamic value pointer and true if found, or nil and false if the scope or key does not exist.
func (s *Store) Find(scope string, key string) (*dynamic.Value, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Find(key)
}

// Push inserts or updates a raw key-value pair within the specified scope.
// If the target scope does not exist, it will be automatically initialized.
// Returns the Store instance to allow method chaining.
func (s *Store) Push(scope string, key string, arg any) *Store {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		ctx = newScope(s.clock)
	}

	s.scopes[scope] = ctx.Push(key,
		newEntry(s.clock, arg),
	)

	return s
}

// RemoveScope completely deletes an entire scope and all its stored arguments.
// Returns true if the scope existed and was removed, or false otherwise.
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

// RemoveArgument deletes a specific key from a scope.
// Returns the removed dynamic value pointer and true if found, or nil and false if missing.
func (s *Store) RemoveArgument(scope, key string) (*dynamic.Value, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ctx, ok := s.scopes[scope]
	if !ok {
		return nil, false
	}

	return ctx.Remove(key)
}

// RetainOnly purges all existing scopes from the store EXCEPT those included in the provided set.
// Useful for garbage collecting scopes when switching navigation routes or screens.
func (s *Store) RetainOnly(scopes set.Set[string]) *Store {
	s.mu.Lock()
	defer s.mu.Unlock()

	for scope := range s.scopes {
		if !scopes.Has(scope) {
			delete(s.scopes, scope)
		}
	}

	return s
}

// Find is a strongly-typed generic helper that retrieves and casts an argument from a scope
// using a typed Key[T]. Returns the zero value of T and false if not found or if casting fails.
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

// Push is a strongly-typed generic helper that writes a value of type T
// into a specified scope using a typed Key[T].
func Push[T any](
	store *Store,
	scope string,
	key Key[T],
	arg T,
) *Store {
	return store.Push(scope, key.Code(), arg)
}

// Update retrieves an existing value of type T, applies an in-place Updater function to it,
// and persists the result. Returns false without calling updater if the key or scope is missing.
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

// Upsert updates an existing value of type T or initializes it with its zero-value first if missing,
// applies the Updater function, and persists the modified state.
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

// Remove is a strongly-typed generic helper that deletes and returns an argument of type T
// associated with a Key[T] from a scope.
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
