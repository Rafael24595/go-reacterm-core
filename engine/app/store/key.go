package store

// Key defines a strongly-typed string key associated with a specific value type T.
// It provides a fluent, ergonomic API for interacting with a Store.
type Key[T any] string

// Type returns the zero-value of type T. Useful for reflection or type inspection.
func (t Key[T]) Type() T {
	var zero T
	return zero
}

// Code returns the underlying string representation of the key.
func (t Key[T]) Code() string {
	return string(t)
}

// Get retrieves the value associated with this key in the specified scope.
// Returns the typed value and true if found, or zero-value and false otherwise.
func (t Key[T]) Get(c *Store, scope string) (T, bool) {
	return Find(c, scope, t)
}

// Set stores the provided argument of type T under this key within the given scope.
// Returns the key to support method chaining.
func (t Key[T]) Set(c *Store, scope string, arg T) Key[T] {
	Push(c, scope, t, arg)
	return t
}

// Update modifies an existing value of type T using the provided updater function.
// Returns the key to support method chaining.
func (t Key[T]) Update(c *Store, scope string, updater Updater[T]) Key[T] {
	Update(c, scope, t, updater)
	return t
}

// Upsert updates an existing value or initializes it with its zero-value before updating.
// Returns the key to support method chaining.
func (t Key[T]) Upsert(c *Store, scope string, updater Updater[T]) Key[T] {
	Upsert(c, scope, t, updater)
	return t
}

// Take retrieves and removes the value associated with this key from the scope.
// Returns the typed value and true if found, or zero-value and false otherwise.
func (t Key[T]) Take(c *Store, scope string) (T, bool) {
	return Remove(c, scope, t)
}

// Delete removes the argument associated with this key from the scope without returning its value.
// Returns the key to support method chaining.
func (t Key[T]) Delete(c *Store, scope string) Key[T] {
	Remove(c, scope, t)
	return t
}
