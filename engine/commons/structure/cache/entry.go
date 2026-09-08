package cache

// DefaultMaxUsed defines the default maximum reference weight for a cache entry.
const DefaultMaxUsed uint8 = 8

// entry represents an internal key-value wrapper with a reference counter (used)
// utilized by eviction algorithms like Clock or LFU.
type entry[T comparable, V any] struct {
	key   T
	value V
	used  uint8
}

func newEntry[T comparable, V any](key T, value V) *entry[T, V] {
	return &entry[T, V]{
		key:   key,
		value: value,
		used:  1,
	}
}

// touch increments the reference count without overflowing uint8, capped at max.
func (e *entry[T, V]) touch(max uint8) {
	inc := e.used
	if e.used < ^uint8(0) {
		inc += 1
	}

	e.used = min(max, inc)
}

// cool decrements the reference count if greater than 0.
// Returns true when the entry reaches zero references (ready for eviction).
func (e *entry[T, V]) cool() bool {
	if e.used == 0 {
		return true
	}

	e.used--
	return false
}
