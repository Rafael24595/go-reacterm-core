package dict

// Pair represents a key-value pair.
type Pair[K any, V any] struct {
	// Key is the key of the pair.
	Key   K
	// Value is the value associated with the key.
	Value V
}

// NewPair creates a new Pair with the given key and value.
func NewPair[K any, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{
		Key:   key,
		Value: value,
	}
}

// P is a shorthand function to create a new Pair.
func P[K any, V any](key K, value V) Pair[K, V] {
	return NewPair(key, value)
}
