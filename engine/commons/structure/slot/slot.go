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

	s.value = &v
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
