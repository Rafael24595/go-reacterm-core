package wrap

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestWordFragMeasure_CacheSameCols(t *testing.T) {
	w := newWordFrag(
		text.NewFrag("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	first := w.measureWith(80, resolver)
	second := w.measureWith(80, resolver)

	assert.Equal(t, first, second)
	assert.Equal(t, winsize.Cols(80), w.cols)

	assert.Equal(t, 1, calls)
}

func TestWordFragMeasure_RecalculateOnColsChange(t *testing.T) {
	w := newWordFrag(
		text.NewFrag("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	_ = w.measureWith(80, resolver)
	m40 := w.measureWith(40, resolver)

	assert.Equal(t, winsize.Cols(40), w.cols)
	assert.Equal(t, m40, w.measure)

	assert.Equal(t, 2, calls)
}

func TestWordFragMeasure_CacheAfterColsChange(t *testing.T) {
	w := newWordFrag(
		text.NewFrag("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	w.measureWith(80, resolver)
	w.measureWith(40, resolver)
	w.measureWith(40, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestWordFragMeasure_RecalculateWhenReturningToPreviousCols(t *testing.T) {
	w := newWordFrag(
		text.NewFrag("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	w.measureWith(80, resolver)
	w.measureWith(40, resolver)
	w.measureWith(80, resolver)

	assert.Equal(t, uint(3), calls)
}
