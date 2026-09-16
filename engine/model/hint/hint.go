package hint

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

type provider[T math.Number] func(max T) T

// Size represents a dynamic or fixed dimensional constraint over a maximum space limit.
type Size[T math.Number] struct {
	provider provider[T]
}

// Fixed constructs a Size that always resolves to a static dimension value.
func Fixed[T math.Number](size T) Size[T] {
	return Size[T]{
		provider: func(T) T {
			return size
		},
	}
}

// Percent constructs a Size that calculates its dimension as a percentage of available space.
func Percent[T math.Number](size T) Size[T] {
	return Size[T]{
		provider: func(max T) T {
			return (size * max) / 100
		},
	}
}

// Maximize constructs a Size that expands to consume all available space.
func Maximize[T math.Number]() Size[T] {
	return Size[T]{
		provider: func(max T) T {
			return max
		},
	}
}

// Min evaluates the calculated size provider and clamps the result to never exceed the max boundary.
func (h Size[T]) Min(max T) T {
	if h.provider == nil {
		assert.Unreachable("the size provider cannot be nil")
		return 0
	}

	return min(
		h.provider(max), max,
	)
}
