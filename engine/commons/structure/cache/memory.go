package cache

// memoryCache provides an unbounded in-memory cache implementation wrapping a Go map.
type memoryCache[T comparable, K any] struct {
	items map[T]K
}

// NewMemory creates a new unbounded in-memory Cache instance.
func NewMemory[T comparable, K any]() Cache[T, K] {
	mem := &memoryCache[T, K]{
		items: make(map[T]K),
	}

	return Cache[T, K]{
		Get: mem.Get,
		Put: mem.Put,
		Del: mem.Del,
		Len: mem.Len,
		Cls: mem.Cls,
	}
}

// Get retrieves the value associated with key from the cache.
// Returns the value and true if present, or zero value and false otherwise.
func (c *memoryCache[T, K]) Get(key T) (K, bool) {
	item, ok := c.items[key]
	return item, ok
}

// Put inserts or overwrites a key-value pair in the cache.
func (c *memoryCache[T, K]) Put(key T, val K) {
	c.items[key] = val
}

// Del removes a key and its value from the cache.
func (c *memoryCache[T, K]) Del(key T) {
	delete(c.items, key)
}

// Len returns the total number of items stored in the cache.
func (c *memoryCache[T, K]) Len() uint {
	return uint(len(c.items))
}

// Cls removes all items from the cache.
func (c *memoryCache[T, K]) Cls() {
	clear(c.items)
}
