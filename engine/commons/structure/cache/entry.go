package cache

const DefaultMaxUsed uint8 = 8

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

func (e *entry[T, V]) touch(max uint8) {
	inc := e.used
	if e.used < ^uint8(0) {
		inc += 1
	}

	e.used = min(max, inc)
}

func (e *entry[T, V]) cool() bool {
	if e.used == 0 {
		return true
	}

	e.used--
	return false
}
