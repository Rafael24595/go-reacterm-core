package cache

type Cache[T, K any] struct {
	Get func(key T) (K, bool)
	Put func(key T, val K)
	Del func(T)
	Len func() uint
	Cls func()
}
