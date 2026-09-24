package winsize_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestWinsize_SubClampZero(t *testing.T) {
	w1 := winsize.New(10, 20)
	w2 := winsize.New(15, 5)

	assert.Equal(t, 0, w1.Rows.Sub(w2.Rows))
	assert.Equal(t, 15, w1.Cols.Sub(w2.Cols))
}

func TestWinsize_Eq(t *testing.T) {
	w1 := winsize.New(10, 10)
	w2 := winsize.New(10, 10)
	w3 := winsize.New(10, 20)

	assert.True(t, w1.Eq(w2))
	assert.False(t, w1.Eq(w3))
}
