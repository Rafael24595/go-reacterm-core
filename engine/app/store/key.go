package store

type Key[T any] string

func (t Key[T]) Type() T {
	var zero T
	return zero
}

func (t Key[T]) Code() string {
	return string(t)
}

func (t Key[T]) Get(c *Store, scope string) (T, bool) {
	return Find(c, scope, t)
}

func (t Key[T]) Set(c *Store, scope string, arg T) Key[T] {
	Push(c, scope, t, arg)
	return t
}

func (t Key[T]) Delete(c *Store, scope string) Key[T] {
	Remove(c, scope, t)
	return t
}
