package dict

import (
	"iter"
	"sync"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/list"
)

const (
	// ErrorImmutableSource indicates that the source is immutable and cannot be modified.
	ErrorImmutableSource = "cannot modify an immutable source"
)

// LinkedMap represents an ordered map that maintains insertion order
// using an internal doubly linked list alongside a hash map for O(1) lookups.
type LinkedMap[K comparable, V any] struct {
	init sync.Once
	inmu bool
	list *list.List[Pair[K, V]]
	data map[K]*list.Item[Pair[K, V]]
}

// NewLinkedMap creates and initializes a new mutable LinkedMap.
func NewLinkedMap[K comparable, V any]() *LinkedMap[K, V] {
	return new(LinkedMap[K, V]).Init()
}

// NewInmutableLinkedMap creates an immutable LinkedMap pre-populated with the given pairs.
// Any subsequent modification attempt (Set, Delete, Merge, etc.) will trigger an assertion panic.
func NewInmutableLinkedMap[K comparable, V any](pairs ...Pair[K, V]) *LinkedMap[K, V] {
	linked := NewLinkedMap[K, V]().Init()
	linked.inmu = true

	for _, p := range pairs {
		linked.set(p)
	}

	return linked
}

func (m *LinkedMap[K, V]) lazyInit() *LinkedMap[K, V] {
	return m.Init()
}

// Init initializes or resets the LinkedMap instance.
func (m *LinkedMap[K, V]) Init() *LinkedMap[K, V] {
	m.init.Do(func() {
		m.list = list.New[Pair[K, V]]()
		m.data = make(map[K]*list.Item[Pair[K, V]])
	})
	return m
}

// Size returns the number of key-value pairs stored in the map.
func (m *LinkedMap[K, V]) Size() uint {
	m.lazyInit()
	return uint(len(m.data))
}

// Exists checks if a given key is present in the map.
func (m *LinkedMap[K, V]) Exists(k K) bool {
	_, exists := m.Get(k)
	return exists
}

// Get retrieves the value associated with key k.
// Returns the value and true if found, or zero value and false if not present.
func (m *LinkedMap[K, V]) Get(k K) (V, bool) {
	m.lazyInit()

	if item, exists := m.data[k]; exists {
		return item.Data.Value, true
	}

	var zero V
	return zero, false
}

// Set inserts or updates a key-value pair.
// Returns the previous value and true if updating an existing key, or zero value and false if inserting new.
func (m *LinkedMap[K, V]) Set(k K, v V) (V, bool) {
	var old V

	if m.inmu {
		assert.Unreachable(ErrorImmutableSource)
		return old, false
	}

	m.lazyInit()

	pair := NewPair(k, v)
	return m.set(pair)
}

// SetPairs inserts or updates multiple key-value pairs.
// Returns the number of newly added keys (excluding updates) and true on success.
func (m *LinkedMap[K, V]) SetPairs(pairs ...Pair[K, V]) (uint, bool) {
	if m.inmu {
		assert.Unreachable(ErrorImmutableSource)
		return 0, false
	}

	if len(pairs) == 0 {
		return 0, true
	}

	m.lazyInit()

	var added uint
	for _, p := range pairs {
		if _, exists := m.set(p); !exists {
			added += 1
		}
	}

	return added, true
}

func (m *LinkedMap[K, V]) set(pair Pair[K, V]) (V, bool) {
	var old V

	item, exists := m.data[pair.Key]
	if !exists {
		m.data[pair.Key] = m.list.Push(pair)
		return old, false
	}

	old = item.Data.Value
	item.Data = pair

	return old, true
}

// Merge inserts all key-value pairs from the other map into m, overwriting existing keys.
// Returns the count of newly inserted keys and true.
func (m *LinkedMap[K, V]) Merge(other *LinkedMap[K, V]) (uint, bool) {
	if m.inmu {
		assert.Unreachable(ErrorImmutableSource)
		return 0, false
	}

	if other == nil {
		return 0, true
	}

	m.lazyInit()

	var added uint
	for p := range other.Pairs() {
		if _, exists := m.set(p); !exists {
			added += 1
		}
	}

	return added, true
}

// Supplement inserts key-value pairs from other into m ONLY if the key does not exist yet.
// Existing keys in m remain untouched. Returns the count of newly added keys and true.
func (m *LinkedMap[K, V]) Supplement(other *LinkedMap[K, V]) (uint, bool) {
	if m.inmu {
		assert.Unreachable(ErrorImmutableSource)
		return 0, false
	}

	if other == nil {
		return 0, true
	}

	m.lazyInit()
	var added uint

	for p := range other.Pairs() {
		if !m.Exists(p.Key) {
			m.set(p)
			added++
		}
	}

	return added, true
}

// Delete removes a key and its associated value from the map.
// Returns the deleted value and true if found, or zero value and false if the key was missing.
func (m *LinkedMap[K, V]) Delete(k K) (V, bool) {
	var old V

	if m.inmu {
		assert.Unreachable(ErrorImmutableSource)
		return old, false
	}

	m.lazyInit()

	item, exists := m.data[k]
	if !exists {
		return old, false
	}

	old = item.Data.Value

	m.list.Delete(item)
	delete(m.data, k)

	return old, true
}

// All returns a two-variable iterator (iter.Seq2) yielding key-value pairs in insertion order.
func (m *LinkedMap[K, V]) All() iter.Seq2[K, V] {
	m.lazyInit()

	return func(yield func(K, V) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Key, item.Data.Value) {
				return
			}
		}
	}
}

// Pairs returns an iterator yielding Pair[K, V] structs in insertion order.
func (m *LinkedMap[K, V]) Pairs() iter.Seq[Pair[K, V]] {
	m.lazyInit()

	return func(yield func(Pair[K, V]) bool) {
		for item := range m.list.All() {
			if !yield(item.Data) {
				return
			}
		}
	}
}

// ToPairsSlice returns a slice containing all Pair[K, V] items in insertion order.
func (m *LinkedMap[K, V]) ToPairsSlice() []Pair[K, V] {
	m.lazyInit()
	pairs := make([]Pair[K, V], 0, len(m.data))
	for p := range m.Pairs() {
		pairs = append(pairs, p)
	}
	return pairs
}

// Keys returns an iterator yielding all keys in insertion order.
func (m *LinkedMap[K, V]) Keys() iter.Seq[K] {
	m.lazyInit()

	return func(yield func(K) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Key) {
				return
			}
		}
	}
}

// ToKeysSlice returns a slice of all keys in insertion order.
func (m *LinkedMap[K, V]) ToKeysSlice() []K {
	m.lazyInit()
	keys := make([]K, 0, len(m.data))
	for k := range m.Keys() {
		keys = append(keys, k)
	}
	return keys
}

// Values returns an iterator yielding all values in insertion order.
func (m *LinkedMap[K, V]) Values() iter.Seq[V] {
	m.lazyInit()

	return func(yield func(V) bool) {
		for item := range m.list.All() {
			if !yield(item.Data.Value) {
				return
			}
		}
	}
}

// ToValuesSlice returns a slice of all values in insertion order.
func (m *LinkedMap[K, V]) ToValuesSlice() []V {
	m.lazyInit()
	values := make([]V, 0, len(m.data))
	for v := range m.Values() {
		values = append(values, v)
	}
	return values
}

// Clone creates a shallow copy of the LinkedMap.
// Optionally takes a boolean flag to set whether the clone should be immutable.
func (m *LinkedMap[K, V]) Clone(inmu ...bool) *LinkedMap[K, V] {
	m.lazyInit()

	immutable := m.inmu
	if len(inmu) > 0 {
		immutable = inmu[0]
	}

	cloned := NewLinkedMap[K, V]()
	cloned.inmu = immutable

	for p := range m.Pairs() {
		cloned.set(p)
	}

	return cloned
}
