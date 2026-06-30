package stack

type Stack[T any] struct {
	data []T
}

func New[T any](limit uint) *Stack[T] {
	return &Stack[T]{
		data: make([]T, 0, limit),
	}
}

func (s *Stack[T]) Len() uint {
	return uint(len(s.data))
}

func (s *Stack[T]) Cap() uint {
	return uint(cap(s.data))
}

func (s *Stack[T]) Items() []T {
	items := make([]T, len(s.data))
	for i := range s.data {
		index := len(s.data) - 1 - i
		items[i] = s.data[index]
	}
	return items
}

func (s *Stack[T]) Peek() (T, bool) {
	var zero T

	if len(s.data) == 0 {
		return zero, false
	}

	return s.data[len(s.data)-1], true
}

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

func (s *Stack[T]) Clear() {
	var zero T

	for i := range s.data {
		s.data[i] = zero
	}

	s.data = s.data[:0]
}
