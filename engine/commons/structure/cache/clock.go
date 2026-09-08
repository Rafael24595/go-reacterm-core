package cache

const defaultClockCapacity = 256

// clockCache implements an eviction policy based on the Second Chance (Clock) algorithm.
// It maintains a circular-like buffer of entries with reference bits to approximate LRU behavior in O(1) time.
type clockCache[T comparable, V any] struct {
	items map[T]int
	slots []*entry[T, V]

	hand int
	len  int

	maxUsed uint8
}

// NewClock initializes a Cache backed by a Clock eviction policy.
// If size is omitted or zero, it defaults to 256 slots.
func NewClock[T comparable, V any](size ...uint) Cache[T, V] {
	c := newClock[T, V](size...)

	return Cache[T, V]{
		Get: c.Get,
		Put: c.Put,
		Del: c.Del,
		Len: c.Len,
		Cls: c.Cls,
	}
}

func newClock[T comparable, V any](size ...uint) *clockCache[T, V] {
	capacity := uint(defaultClockCapacity)
	if len(size) > 0 && size[0] > 0 {
		capacity = size[0]
	}

	return &clockCache[T, V]{
		items:   make(map[T]int, capacity),
		slots:   make([]*entry[T, V], capacity),
		maxUsed: DefaultMaxUsed,
	}
}

// Get retrieves a key's value from the cache and updates its reference score (touch).
// Returns the value and true if found, or zero value and false otherwise.
func (c *clockCache[T, V]) Get(key T) (V, bool) {
	if idx, ok := c.items[key]; ok {
		e := c.slots[idx]
		e.touch(c.maxUsed)
		return e.value, true
	}

	var zero V
	return zero, false
}

// Put inserts or updates a key-value pair in the cache.
// If the cache is full, it triggers a Clock eviction to free space for the new entry.
func (c *clockCache[T, V]) Put(key T, value V) {
	if idx, ok := c.items[key]; ok {
		entry := c.slots[idx]

		entry.value = value
		entry.touch(c.maxUsed)

		return
	}

	entry := newEntry(key, value)

	if c.len < len(c.slots) {
		c.slots[c.len] = entry
		c.items[key] = c.len

		c.len++
		return
	}

	c.evict()

	idx := c.hand

	c.slots[idx] = entry
	c.items[key] = idx

	c.incHand()
}

// Del removes a key from the cache in O(1) time by swapping the target slot with the last active slot.
func (c *clockCache[T, V]) Del(key T) {
	idx, ok := c.items[key]
	if !ok {
		return
	}

	delete(c.items, key)

	lastIdx := c.len - 1
	if idx != lastIdx {
		lastEntry := c.slots[lastIdx]

		c.slots[idx] = lastEntry
		c.items[lastEntry.key] = idx
	}

	c.slots[lastIdx] = nil
	c.len -= 1

	if c.len > 0 && c.hand >= c.len {
		c.hand = 0
	}
}

// Len returns the current number of cached elements.
func (c *clockCache[T, V]) Len() uint {
	return uint(c.len)
}

// Cls clears all items from the cache and resets the hand pointer.
func (c *clockCache[T, V]) Cls() {
	clear(c.items)

	for i := 0; i < c.len; i++ {
		c.slots[i] = nil
	}

	c.len = 0
	c.hand = 0
}

func (c *clockCache[T, V]) evict() {
	for !c.slots[c.hand].cool() {
		c.incHand()
	}

	entry := c.slots[c.hand]
	delete(c.items, entry.key)
}

func (c *clockCache[T, V]) incHand() {
	c.hand += 1
	if c.hand >= len(c.slots) {
		c.hand = 0
	}
}
