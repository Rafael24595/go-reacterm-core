package slot

import "sync"

// Slot is a thread-safe container holding at most one value at a time.
type Slot[T any] struct {
	mu    sync.Mutex
	value *T
}

// Slot is a thread-safe container holding at most one value at a time.
func New[T any]() *Slot[T] {
	return &Slot[T]{}
}

// Set replaces or populates the value inside the slot in a thread-safe manner.
func (s *Slot[T]) Set(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val := v
	s.value = &val
}

// Take atomically retrieves and clears the value held inside the slot.
// Returns false if the slot was empty.
func (s *Slot[T]) Take() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.value == nil {
		var zero T
		return zero, false
	}

	v := *s.value
	s.value = nil

	return v, true
}

func (s *Slot[T]) Peek() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.value == nil {
		var zero T
		return zero, false
	}

	return *s.value, true
}
