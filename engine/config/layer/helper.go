package layer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

// Fixed sets a fixed-size chunk strategy for the layer.
func Fixed[T math.Number](chk T) Option[T] {
	return WithChunk(chunk.Fixed(chk))
}

// Percent sets a percentage-based chunk strategy for the layer.
func Percent[T math.Number](chk T) Option[T] {
	return WithChunk(chunk.Percent(chk))
}
