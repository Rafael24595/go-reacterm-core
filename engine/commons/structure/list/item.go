package list

// Item represents a node in a doubly linked list.
type Item[T any] struct {
	// next points to the next item in the list. It is nil if the item is not part of any list.
	next *Item[T]
	// prev points to the previous item in the list. It is nil if the item is not part of any list.
	prev *Item[T]
	// list points to the list that contains this item. It is nil if the item is not part of any list.
	list *List[T]
	// Data holds the value of the list item.
	Data T
}

func newItem[T any](data T) *Item[T] {
	return &Item[T]{
		Data: data,
	}
}

// Next returns the next list item and true, or (nil, false) if e is the last item.
func (e *Item[T]) Next() (*Item[T], bool) {
	if e.list == nil {
		return nil, false
	}

	next := e.next
	if next == &e.list.root {
		return nil, false
	}

	return next, true
}

// Prev returns the previous list item and true, or (nil, false) if e is the first item.
func (e *Item[T]) Prev() (*Item[T], bool) {
	if e.list == nil {
		return nil, false
	}

	prev := e.prev
	if prev == &e.list.root {
		return nil, false
	}

	return prev, true
}
