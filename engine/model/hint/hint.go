package hint

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type provider[T math.Number] func(max T) T

type Size[T math.Number] struct {
	provider provider[T]
}

func Fixed[T math.Number](size T) Size[T] {
	return Size[T]{
		provider: func(T) T {
			return size
		},
	}
}

func Percent[T math.Number](size T) Size[T] {
	return Size[T]{
		provider: func(max T) T {
			return (size * max) / 100
		},
	}
}

func Maximize[T math.Number]() Size[T] {
	return Size[T]{
		provider: func(max T) T {
			return max
		},
	}
}

func (h Size[T]) Min(max T) T {
	return min(
		h.provider(max), max,
	)
}
