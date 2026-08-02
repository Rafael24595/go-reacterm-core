package cache

const defaultClockCapacity = 256

type clockCache[T comparable, V any] struct {
	items map[T]int
	slots []*entry[T, V]

	hand int
	len  int

	maxUsed uint8
}

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

func (c *clockCache[T, V]) Get(key T) (V, bool) {
	if idx, ok := c.items[key]; ok {
		e := c.slots[idx]
		e.touch(c.maxUsed)
		return e.value, true
	}

	var zero V
	return zero, false
}

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

func (c *clockCache[T, V]) Len() uint {
	return uint(c.len)
}

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
