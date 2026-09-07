package list

import (
	"iter"
	"sync"
)

// List represents a doubly linked list.
// The zero value for List is empty and ready to use.
type List[T any] struct {
	init sync.Once
	root Item[T]
	size uint
}

// New creates and initializes a new empty List.
func New[T any]() *List[T] {
	return new(List[T]).Init()
}

func (l *List[T]) lazyInit() *List[T] {
	return l.Init()
}

// Init initializes or clears list l.
func (l *List[T]) Init() *List[T] {
	l.init.Do(func() {
		l.root.next = &l.root
		l.root.prev = &l.root
		l.size = 0
	})
	return l
}

// Size returns the number of elements in list l.
func (l *List[T]) Size() uint {
	l.lazyInit()

	return l.size
}

// First returns the first item of list l and true, or (nil, false) if the list is empty.
func (l *List[T]) First() (*Item[T], bool) {
	l.lazyInit()

	if l.size == 0 {
		return nil, false
	}

	return l.root.next, true
}

// Last returns the last item of list l and true, or (nil, false) if the list is empty.
func (l *List[T]) Last() (*Item[T], bool) {
	l.lazyInit()

	if l.size == 0 {
		return nil, false
	}

	return l.root.prev, true
}

// Unshift inserts a new element at the front of list l and returns the created Item.
func (l *List[T]) Unshift(data T) *Item[T] {
	l.lazyInit()

	item := newItem(data)
	return l.insert(item, &l.root)
}

// Push inserts a new element at the back of list l and returns the created Item.
func (l *List[T]) Push(data T) *Item[T] {
	l.lazyInit()

	item := newItem(data)
	return l.insert(item, l.root.prev)
}

func (l *List[T]) insert(it, at *Item[T]) *Item[T] {
	it.prev = at
	it.next = at.next

	it.prev.next = it
	it.next.prev = it

	it.list = l
	l.size += 1

	return it
}

// Delete removes element e from list l if e belongs to l.
// It returns the deleted value and true if successful, or zero value and false otherwise.
func (l *List[T]) Delete(e *Item[T]) (T, bool) {
	if e == nil || e.list != l {
		var zero T
		return zero, false
	}

	deleted := e.Data

	e.prev.next = e.next
	e.next.prev = e.prev

	var zero T
	e.Data = zero

	e.next = nil
	e.prev = nil
	e.list = nil

	if l.size > 0 {
		l.size -= 1
	}

	return deleted, true
}

// All returns an iterator (iter.Seq) yielding each *Item in the list from first to last.
func (l *List[T]) All() iter.Seq[*Item[T]] {
	return func(yield func(*Item[T]) bool) {
		l.lazyInit()

		for i := l.root.next; i != &l.root; {
			next := i.next
			if !yield(i) {
				return
			}
			if i.list != nil {
				i = i.next
			} else {
				i = next
			}
		}
	}
}
