package chunk

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
)

const maxChunk = 100

const (
	// ErrorChunkSize defines the panic message when a percentage chunk exceeds 100%.
	ErrorChunkSize = "chunk value should be less or equals than %d"
)

// chunkAdapter defines a function type that takes a size of type T and returns an adapted size of the same type.
type chunkAdapter[T math.Number] func(size T) T

// Chunk represents a dimension adapter configuration.
type Chunk[T math.Number] struct {
	isAnemic bool
	Adapter  chunkAdapter[T]
	Sized    bool
}

// Dynamic creates an unconstrained (anemic) Chunk that adapts dynamically to available space.
func Dynamic[T math.Number]() Chunk[T] {
	return Chunk[T]{
		isAnemic: true,
		Adapter:  fixAdapter(T(0)),
		Sized:    false,
	}
}

// Fixed creates a Chunk with a fixed absolute size T, clamped to maximum available space.
func Fixed[T math.Number](fix T) Chunk[T] {
	return Chunk[T]{
		isAnemic: false,
		Adapter:  fixAdapter(fix),
		Sized:    true,
	}
}

// Percent creates a Chunk representing a percentage (0 to 100) of the available space.
// Panics if chunk exceeds 100.
func Percent[T math.Number](chunk T) Chunk[T] {
	if chunk > maxChunk {
		assert.Unreachable(ErrorChunkSize, maxChunk)
		chunk = maxChunk
	}

	return Chunk[T]{
		isAnemic: false,
		Adapter:  perAdapter(chunk),
		Sized:    true,
	}
}

// IsAnemic returns true if the chunk has no fixed or percentage constraints.
func (c Chunk[T]) IsAnemic() bool {
	return c.isAnemic
}

func perAdapter[T math.Number](chunk T) chunkAdapter[T] {
	return func(size T) T {
		return (size * chunk) / 100
	}
}

func fixAdapter[T math.Number](chunk T) chunkAdapter[T] {
	return func(size T) T {
		if chunk > size {
			return size
		}
		return chunk
	}
}
