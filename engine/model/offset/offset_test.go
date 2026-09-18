package offset_test

import (
	"math"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

func TestOffset_Add(t *testing.T) {
	a := offset.Offset(10)
	b := offset.Offset(20)

	got := a.Add(b)
	assert.Equal(t, got, offset.Offset(30))

	maxVal := offset.Offset(math.MaxUint32)
	gotMax := maxVal.Add(offset.Offset(100))
	assert.Equal(t, gotMax, offset.Offset(math.MaxUint32))
}

func TestOffset_Sub(t *testing.T) {
	a := offset.Offset(30)
	b := offset.Offset(10)

	got := a.Sub(b)
	assert.Equal(t, got, offset.Offset(20))

	gotZero := b.Sub(a)
	assert.Equal(t, gotZero, offset.Offset(0))
}