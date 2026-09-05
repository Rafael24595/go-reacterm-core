package heap

import "cmp"

// Heap represents a binary heap of elements ordered by a custom comparator function.
type Heap[T any] struct {
	items []T
	less  func(a, b T) bool
}

// New creates an empty Heap with a custom comparator function.
func New[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less:  less,
	}
}

// NewMin creates an empty Min-Heap for ordered types.
func NewMin[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less: func(a, b T) bool {
			return a < b
		},
	}
}

// NewMax creates an empty Max-Heap for ordered types.
func NewMax[T cmp.Ordered]() *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less: func(a, b T) bool {
			return a > b
		},
	}
}

// NewMinBy creates a Min-Heap ordered by a extracted key.
func NewMinBy[T any, K cmp.Ordered](get func(i T) K) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less: func(a, b T) bool {
			return get(a) < get(b)
		},
	}
}

// NewMaxBy creates a Max-Heap ordered by a extracted key.
func NewMaxBy[T any, K cmp.Ordered](get func(i T) K) *Heap[T] {
	return &Heap[T]{
		items: []T{},
		less: func(a, b T) bool {
			return get(a) > get(b)
		},
	}
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.items)
}

// IsEmpty reports whether the heap contains no items.
func (h *Heap[T]) IsEmpty() bool {
	return len(h.items) == 0
}

// Push adds an item to the heap and restores its invariant.
func (h *Heap[T]) Push(x T) {
	h.items = append(h.items, x)
	h.up(h.Len() - 1)
}

// Pop removes and returns the root item from the heap.
func (h *Heap[T]) Pop() (T, bool) {
	if h.Len() == 0 {
		var zero T
		return zero, false
	}

	top := h.items[0]
	last := h.pop()

	if h.Len() > 0 {
		h.items[0] = last
		h.down(0)
	}

	return top, true
}

// Peek returns the root element without removing it.
func (h *Heap[T]) Peek() (T, bool) {
	if h.Len() == 0 {
		var zero T
		return zero, false
	}
	return h.items[0], true
}

// Clear removes all elements from the heap while maintaining allocated capacity.
func (h *Heap[T]) Clear() {
	clear(h.items)
	h.items = h.items[:0]
}

func (h *Heap[T]) pop() T {
	last := len(h.items) - 1
	item := h.items[last]

	var zero T
	h.items[last] = zero
	h.items = h.items[:last]

	return item
}

func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *Heap[T]) up(index int) {
	for {
		parent := (index - 1) / 2
		if parent == index || !h.less(h.items[index], h.items[parent]) {
			break
		}

		h.swap(parent, index)
		index = parent
	}
}

func (h *Heap[T]) down(index int) bool {
	size := h.Len()

	cursor := index
	for {
		childLeft := 2*cursor + 1
		if childLeft >= size || childLeft < 0 {
			break
		}

		priorityChild := childLeft

		candidateFix := childLeft + 1
		if candidateFix < size && h.less(h.items[candidateFix], h.items[childLeft]) {
			priorityChild = candidateFix
		}

		if !h.less(h.items[priorityChild], h.items[cursor]) {
			break
		}

		h.swap(cursor, priorityChild)
		cursor = priorityChild
	}

	return cursor > index
}
