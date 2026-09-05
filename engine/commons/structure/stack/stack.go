package stack

// Stack represents a generic bounded stack container.
type Stack[T any] struct {
	data []T
}

// New initializes a new Stack with a fixed maximum capacity limit.
func New[T any](limit uint) *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0, limit),
	}
}

// Len returns the current number of elements stored in the stack.
func (s *Stack[T]) Len() uint {
	return uint(len(s.data))
}

// Cap returns the maximum capacity of the stack.
func (s *Stack[T]) Cap() uint {
	return uint(cap(s.data))
}

// Items returns a fresh slice with all items in LIFO order (top element first).
func (s *Stack[T]) Items() []T {
	n := len(s.data)
	items := make([]T, n)
	for i := range n {
		items[i] = s.data[n-1-i]
	}
	return items
}

// Peek returns the top element without removing it from the stack.
func (s *Stack[T]) Peek() (T, bool) {
	var zero T

	if len(s.data) == 0 {
		return zero, false
	}

	return s.data[len(s.data)-1], true
}

// Push adds an item to the top of the stack.
// If the stack is at full capacity, the oldest item (bottom) is evicted and returned as discarded.
func (s *Stack[T]) Push(item T) (T, bool) {
	if len(s.data) < cap(s.data) {
		var zero T
		s.data = append(s.data, item)
		return zero, false
	}

	discarded := s.data[0]

	copy(s.data, s.data[1:])
	s.data[len(s.data)-1] = item

	return discarded, true
}

// Pop removes and returns the top element from the stack.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.data) == 0 {
		return zero, false
	}

	last := len(s.data) - 1
	node := s.data[last]

	s.data[last] = zero

	s.data = s.data[:last]

	return node, true
}

// Clear removes all elements and zeroes memory to prevent leaks.
func (s *Stack[T]) Clear() {
	clear(s.data)
	s.data = s.data[:0]
}
