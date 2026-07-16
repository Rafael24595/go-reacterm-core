package wrap

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestWordFragMeasure_CacheSameCols(t *testing.T) {
	frg := frag.FromString("golang")
	wrd := newWordFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	first := wrd.measureWith(80, resolver)
	second := wrd.measureWith(80, resolver)

	assert.Equal(t, first, second)
	assert.Equal(t, winsize.Cols(80), wrd.cols)

	assert.Equal(t, 1, calls)
}

func TestWordFragMeasure_RecalculateOnColsChange(t *testing.T) {
	frg := frag.FromString("golang")
	wrd := newWordFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	_ = wrd.measureWith(80, resolver)
	m40 := wrd.measureWith(40, resolver)

	assert.Equal(t, winsize.Cols(40), wrd.cols)
	assert.Equal(t, m40, wrd.measure)

	assert.Equal(t, 2, calls)
}

func TestWordFragMeasure_CacheAfterColsChange(t *testing.T) {
	frg := frag.FromString("golang")
	wrd := newWordFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	wrd.measureWith(80, resolver)
	wrd.measureWith(40, resolver)
	wrd.measureWith(40, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestWordFragMeasure_RecalculateWhenReturningToPreviousCols(t *testing.T) {
	frg := frag.FromString("golang")
	wrd := newWordFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	wrd.measureWith(80, resolver)
	wrd.measureWith(40, resolver)
	wrd.measureWith(80, resolver)

	assert.Equal(t, uint(3), calls)
}
