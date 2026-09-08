package layer

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/config/chunk"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
)

type config[T math.Number] struct {
	unit   drawable.Unit
	chunk  chunk.Chunk[T]
	static bool
}

// Layer represents a rendered layout unit with its sizing policy and value.
type Layer[T math.Number] struct {
	config config[T]
	Value  T
	Status bool
}

// New constructs a new Layer instance wrapping a drawable.Unit and options.
func New[T math.Number](unit drawable.Unit, opts ...Option[T]) Layer[T] {
	return FromLayer(
		defaultConfig[T](unit), opts...,
	)
}

// FromLayer copies an existing Layer base and applies new Option modifications.
func FromLayer[T math.Number](other Layer[T], opts ...Option[T]) Layer[T] {
	cfg := other
	for _, opt := range opts {
		opt(&cfg)
	}

	unit := cfg.Unit()

	assert.LazyFalse(func() bool {
		return drawable.IsZeroUnit(unit)
	}, "unit '%s' is not defined", unit.Name)

	cfg.Status = true

	return cfg
}

// IsAnemic returns true if the layer has dynamic sizing and is not static.
func (l Layer[T]) IsAnemic() bool {
	return l.Chunk().IsAnemic() && !l.config.static
}

// Unit retrieves the underlying drawable.Unit.
func (l Layer[T]) Unit() drawable.Unit {
	return l.config.unit
}

// Chunk returns the current sizing chunk strategy.
func (l Layer[T]) Chunk() chunk.Chunk[T] {
	return l.config.chunk
}

// Static returns true if the layer dimension is marked as static.
func (l Layer[T]) Static() bool {
	return l.config.static
}
