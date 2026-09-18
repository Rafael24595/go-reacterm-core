package offset

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

// Offset represents a non-negative unsigned integer position or index.
type Offset uint32

// Add performs addition with another Offset, clamping to math.MaxUint32 on overflow.
func (o Offset) Add(of Offset) Offset {
	return math.AddClampMax(o, of)
}

// Sub performs subtraction with another Offset, clamping to zero on underflow.
func (o Offset) Sub(of Offset) Offset {
	return math.SubClampZero(o, of)
}
