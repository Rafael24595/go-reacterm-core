package set

import "maps"

// Set represents a collection of unique comparable items.
type Set[T comparable] map[T]struct{}

// New creates an empty Set with optional pre-allocated capacity.
func New[T comparable](size ...int) Set[T] {
	s := 0

	if len(size) > 0 {
		s = size[0]
	}

	return make(Set[T], s)
}

// From constructs a Set pre-populated with items from a slice.
func From[T comparable](slice ...T) Set[T] {
	s := New[T](len(slice))

	for _, v := range slice {
		s.Add(v)
	}

	return s
}

// Add inserts one or multiple items into the set.
func (s Set[T]) Add(v ...T) {
	for _, t := range v {
		s[t] = struct{}{}
	}
}

// Remove deletes one or multiple items from the set.
func (s Set[T]) Remove(v ...T) {
	for _, t := range v {
		delete(s, t)
	}
}

// Merge adds all elements from another set into the current set.
func (s Set[T]) Merge(v Set[T]) {
	for k := range v {
		s.Add(k)
	}
}

// Has reports whether an item exists in the set.
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Any returns true if the set shares at least one common element with another set.
// Operates in O(min(|A|, |B|)) by iterating over the smaller set.
func (s Set[T]) Any(other Set[T]) bool {
	a, b := s, other
	if len(a) > len(b) {
		a, b = b, a
	}

	for k := range a {
		if _, ok := b[k]; ok {
			return true
		}
	}

	return false
}

func (s Set[T]) Slice() []T {
	items := make([]T, 0, len(s))
	for k := range s {
		items = append(items, k)
	}
	return items
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}
