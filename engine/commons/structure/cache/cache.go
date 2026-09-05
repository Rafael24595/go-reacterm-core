package cache

// Cache represents a generic cache structure with basic operations.
type Cache[T, K any] struct {
	// Get retrieves the value associated with the given key.
	Get func(key T) (K, bool)
	// Put adds or updates the value associated with the given key.
	Put func(key T, val K)
	// Del removes the value associated with the given key.
	Del func(key T)
	// Len returns the number of items in the cache.
	Len func() uint
	// Cls clears all items from the cache.
	Cls func()
}
