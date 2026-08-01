package layout

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestWordFragMeasure_CacheSameCols(t *testing.T) {
	frg := frag.FromString("golang")
	wrd := NewFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
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
	wrd := NewFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
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
	wrd := NewFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
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
	wrd := NewFrag(&frg)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return 42
	}

	wrd.measureWith(80, resolver)
	wrd.measureWith(40, resolver)
	wrd.measureWith(80, resolver)

	assert.Equal(t, uint(3), calls)
}

func TestSplitFragAt(t *testing.T) {
	frg := frag.FromString("ziglang")
	wrd := NewFrag(&frg)

	left, right := splitFragAt(wrd, 3)

	assert.NotNil(t, left)
	assert.NotNil(t, right)

	assert.NotSame(t, wrd, left)
	assert.NotSame(t, wrd, right)

	assert.NotSame(t, wrd.Base, left.Base)
	assert.NotSame(t, wrd.Base, right.Base)

	assert.Equal(t, "zig", left.Base.Text())
	assert.Equal(t, "lang", right.Base.Text())
}

func TestSplitFragAt_StartOfFrag(t *testing.T) {
	frg := frag.FromString("abcdef")
	wrd := NewFrag(&frg)

	left, right := splitFragAt(wrd, 0)

	assert.NotNil(t, left)
	assert.NotNil(t, right)

	assert.NotSame(t, wrd, left)
	assert.NotSame(t, wrd, right)

	assert.NotSame(t, wrd.Base, left.Base)
	assert.NotSame(t, wrd.Base, right.Base)

	assert.Equal(t, "abcdef", right.Base.Text())
}

func TestSplitFragAt_EndOfFrag(t *testing.T) {
	frg := frag.FromString("abcdef")
	wrd := NewFrag(&frg)

	left, right := splitFragAt(wrd, 6)

	assert.NotNil(t, left)
	assert.Nil(t, right)

	assert.NotSame(t, wrd, left)

	assert.NotSame(t, wrd.Base, left.Base)

	assert.Equal(t, "abcdef", left.Base.Text())
}

