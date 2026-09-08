package layer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
)

// Option defines a functional configuration option for Layer.
type Option[T math.Number] func(*Layer[T])

func defaultConfig[T math.Number](unit drawable.Unit) Layer[T] {
	config := config[T]{
		unit:   unit,
		chunk:  chunk.Dynamic[T](),
		static: false,
	}

	return Layer[T]{
		config: config,
		Value:  0,
		Status: true,
	}
}

// WithChunk sets a custom sizing strategy.
func WithChunk[T math.Number](chunk chunk.Chunk[T]) Option[T] {
	return func(cfg *Layer[T]) {
		cfg.config.chunk = chunk
	}
}

// WithValue assigns an explicit numerical value to the layer.
func WithValue[T math.Number](value T) Option[T] {
	return func(cfg *Layer[T]) {
		cfg.Value = value
	}
}

// Static marks the layer as static.
func Static[T math.Number]() Option[T] {
	return func(cfg *Layer[T]) {
		cfg.config.static = true
	}
}
