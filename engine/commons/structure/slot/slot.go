package slot

import "sync"

type Slot[T any] struct {
	mu    sync.Mutex
	value *T
}

func New[T any]() *Slot[T] {
	return &Slot[T]{}
}

func (s *Slot[T]) Set(v T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val := v
	s.value = &val
}

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
