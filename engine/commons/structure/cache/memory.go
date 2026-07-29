package cache

type memoryCache[T comparable, K any] struct {
	items map[T]K
}

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

func (c *memoryCache[T, K]) Get(key T) (K, bool) {
	item, ok := c.items[key]
	return item, ok
}

func (c *memoryCache[T, K]) Put(key T, val K) {
	c.items[key] = val
}

func (c *memoryCache[T, K]) Del(key T) {
	delete(c.items, key)
}

func (c *memoryCache[T, K]) Len() uint {
	return uint(len(c.items))
}

func (c *memoryCache[T, K]) Cls() {
	c.items = make(map[T]K)
}
