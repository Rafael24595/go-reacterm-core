package structure_test

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/cache"

type MockCache[T, K any] struct {
	GetCalls uint
	Get      func(key T) (K, bool)
	PutCalls uint
	Put      func(key T, val K)
	DelCalls uint
	Del      func(key T)
	LenCalls uint
	Len      func() uint
	ClsCalls uint
	Cls      func()
}

func NewMockCache[T, K any]() cache.Cache[T, K] {
	mock := MockCache[T, K]{}
	return mock.ToCache()
}

func MockFromCache[T, K any](cache cache.Cache[T, K]) *MockCache[T, K] {
	mock := MockCache[T, K]{}

	mock.Get = func(key T) (K, bool) {
		mock.GetCalls += 1
		return cache.Get(key)
	}

	mock.Put = func(key T, val K) {
		mock.PutCalls += 1
		cache.Put(key, val)
	}

	mock.Del = func(key T) {
		mock.DelCalls += 1
		cache.Del(key)
	}

	mock.Len = func() uint {
		mock.LenCalls += 1
		return cache.Len()
	}

	mock.Cls = func() {
		mock.ClsCalls += 1
		cache.Cls()
	}

	return &mock
}

func (m *MockCache[T, K]) ToCache() cache.Cache[T, K] {
	if m.Get == nil {
		m.Get = func(key T) (K, bool) {
			var zero K
			return zero, false
		}
	}

	if m.Put == nil {
		m.Put = func(key T, val K) {}
	}

	if m.Del == nil {
		m.Del = func(key T) {}
	}

	if m.Len == nil {
		m.Len = func() uint {
			return 0
		}
	}

	if m.Cls == nil {
		m.Cls = func() {}
	}

	return cache.Cache[T, K]{
		Get: m.Get,
		Put: m.Put,
		Del: m.Del,
		Len: m.Len,
		Cls: m.Cls,
	}
}
