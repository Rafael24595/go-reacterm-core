package layout

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func sourceLayout(source line.Line) *Line {
	return &Line{
		Source: source,
		words:  make([]Word, 0),
		frags:  make([]Frag, 0),
	}
}

func emptyLayout() *Line {
	return sourceLayout(line.Line{})
}

func fragsToString(frags []Frag) string {
	var b strings.Builder
	for _, f := range frags {
		b.WriteString(f.Base.Text())
	}
	return b.String()
}


func TestLayoutFindFrags(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("foo"),
			frag.FromString("bar"),
		).
		PushFrags(
			frag.FromString("baz"),
		)

	frags := lne.FindFrags(0)

	assert.Size(t, 2, frags)
	assert.Equal(t, "foo", frags[0].Base.Text())
	assert.Equal(t, "bar", frags[1].Base.Text())

	last := lne.FindFrags(1)

	assert.Size(t, 1, last)
	assert.Equal(t, "baz", last[0].Base.Text())
}

func TestLayoutPushFrags(t *testing.T) {
	lne := emptyLayout()

	lne.PushFrags(frag.FromString("a")).
		PushFrags(
			frag.FromString("b"),
			frag.FromString("c"),
		)

	assert.Size(t, 2, lne.words)
	assert.Size(t, 3, lne.frags)

	assert.Equal(t, uint32(0), lne.words[0].start)
	assert.Equal(t, uint32(1), lne.words[0].end)

	assert.Equal(t, uint32(1), lne.words[1].start)
	assert.Equal(t, uint32(3), lne.words[1].end)
}

func TestLayoutHasAtom(t *testing.T) {
	frg := frag.TextAtom("foo", atom.Wrap)

	lne := emptyLayout().
		PushFrags(frg)

	assert.True(t, lne.HasAtom(0, atom.Wrap))
	assert.False(t, lne.HasAtom(0, atom.Focus))
}

func TestSplitWord(t *testing.T) {
	tests := []struct {
		name            string
		layout          *Line
		cols            winsize.Cols
		remaining       winsize.Cols
		expectedCurrent string
		expectedRest    string
	}{
		{
			name: "word fits completely",
			layout: emptyLayout().PushFrags(
				frag.FromString("golang"),
			),
			cols:            20,
			remaining:       20,
			expectedCurrent: "golang",
			expectedRest:    "",
		},
		{
			name: "split single frag word",
			layout: emptyLayout().PushFrags(
				frag.FromString("ziglang"),
			),
			cols:            4,
			remaining:       4,
			expectedCurrent: "zigl",
			expectedRest:    "ang",
		},
		{
			name: "split fragmented word",
			layout: emptyLayout().PushFrags(
				frag.FromString("go"),
				frag.FromString("la"),
				frag.FromString("ng"),
			),
			cols:            2,
			remaining:       4,
			expectedCurrent: "gola",
			expectedRest:    "ng",
		},
		{
			name: "zero remaining",
			layout: emptyLayout().PushFrags(
				frag.FromString("rust"),
			),
			cols:            5,
			remaining:       0,
			expectedCurrent: "",
			expectedRest:    "rust",
		},
		{
			name: "split inside second frag",
			layout: emptyLayout().PushFrags(
				frag.FromString("cl"),
				frag.FromString("oju"),
				frag.FromString("re"),
			),
			cols:            3,
			remaining:       3,
			expectedCurrent: "clo",
			expectedRest:    "jure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := uint(0)

			ok := tt.layout.SplitWord(
				idx, tt.cols, tt.remaining,
			)

			if tt.expectedCurrent != "" {
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCurrent, fragsToString(tt.layout.FindFrags(idx)))
			}
		})
	}
}

func TestSplitWord_FitsWithoutMutatingLayout(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("golang"),
		)

	ok := lne.SplitWord(0, 80, 80)

	assert.True(t, ok)

	assert.Size(t, 1, lne.words)
	assert.Size(t, 1, lne.frags)

	assert.Equal(t, "golang", fragsToString(lne.FindFrags(0)))

	assert.Equal(t, uint32(0), lne.words[0].start)
	assert.Equal(t, uint32(1), lne.words[0].end)
}

func TestSplitWord_SplitLastFrag(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("go"),
			frag.FromString("la"),
			frag.FromString("ng"),
		)

	ok := lne.SplitWord(0, 80, 5)

	assert.True(t, ok)

	assert.Size(t, 2, lne.words)
	assert.Size(t, 4, lne.frags)

	assert.Equal(t, "golan", fragsToString(lne.FindFrags(0)))
	assert.Equal(t, "n", lne.frags[2].Base.Text())
	assert.Equal(t, "g", fragsToString(lne.FindFrags(1)))
}

func TestLayoutSplitFrag(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("abcdef"),
		)

	lne.splitFrag(0, 0, 3)

	assert.Size(t, 2, lne.frags)
	assert.Size(t, 2, lne.words)

	assert.Equal(t, "abc", lne.frags[0].Base.Text())
	assert.Equal(t, "def", lne.frags[1].Base.Text())

	assert.Equal(t, uint32(0), lne.words[0].start)
	assert.Equal(t, uint32(1), lne.words[0].end)

	assert.Equal(t, uint32(1), lne.words[1].start)
	assert.Equal(t, uint32(2), lne.words[1].end)
}

func TestLayoutSplitFrag_InvalidatesMeasureCache(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("abcdef"),
		)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return winsize.Cols(len(frags))
	}

	lne.measureWith(0, 80, resolver)
	assert.Equal(t, uint(1), calls)

	lne.splitFrag(0, 0, 3)

	lne.measureWith(0, 80, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestLayoutSplitFrag_ShiftsFollowingWords(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("abcdef"),
		).
		PushFrags(
			frag.FromString("foo"),
		).
		PushFrags(
			frag.FromString("bar"),
		)

	lne.splitFrag(0, 0, 3)

	assert.Size(t, 4, lne.words)
	assert.Size(t, 4, lne.frags)

	assert.Equal(t, "abc", fragsToString(lne.FindFrags(0)))
	assert.Equal(t, "def", fragsToString(lne.FindFrags(1)))
	assert.Equal(t, "foo", fragsToString(lne.FindFrags(2)))
	assert.Equal(t, "bar", fragsToString(lne.FindFrags(3)))

	assert.Equal(t, 0, lne.words[0].start)
	assert.Equal(t, 1, lne.words[0].end)

	assert.Equal(t, 1, lne.words[1].start)
	assert.Equal(t, 2, lne.words[1].end)

	assert.Equal(t, 2, lne.words[2].start)
	assert.Equal(t, 3, lne.words[2].end)

	assert.Equal(t, 3, lne.words[3].start)
	assert.Equal(t, 4, lne.words[3].end)
}

func TestLayoutSplitFrag_NoSplit(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("abc"),
		)

	lne.splitFrag(0, 0, 3)

	assert.Size(t, 1, lne.words)
	assert.Size(t, 1, lne.frags)

	assert.Equal(t, "abc", fragsToString(lne.FindFrags(0)))
}

func TestLayoutWordMeasure_CacheSameCols(t *testing.T) {
	lne := emptyLayout().PushFrags(
		frag.FromString("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return 42
	}

	first := lne.measureWith(0, 80, resolver)
	second := lne.measureWith(0, 80, resolver)

	assert.Equal(t, first, second)
	assert.Equal(t, winsize.Cols(80), lne.words[0].cols)

	assert.Equal(t, 1, calls)
}

func TestLayoutWordMeasure_RecalculateOnColsChange(t *testing.T) {
	lne := emptyLayout().PushFrags(
		frag.FromString("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return 42
	}

	_ = lne.measureWith(0, 80, resolver)
	m40 := lne.measureWith(0, 40, resolver)

	assert.Equal(t, winsize.Cols(40), lne.words[0].cols)
	assert.Equal(t, m40, lne.words[0].measure)

	assert.Equal(t, 2, calls)
}

func TestLayoutWordMeasure_CacheAfterColsChange(t *testing.T) {
	lne := emptyLayout().PushFrags(
		frag.FromString("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return 42
	}

	lne.measureWith(0, 80, resolver)
	lne.measureWith(0, 40, resolver)
	lne.measureWith(0, 40, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestLayoutWordMeasure_RecalculateWhenReturningToPreviousCols(t *testing.T) {
	lne := emptyLayout().PushFrags(
		frag.FromString("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...Frag) winsize.Cols {
		calls++
		return 42
	}

	lne.measureWith(0, 80, resolver)
	lne.measureWith(0, 40, resolver)
	lne.measureWith(0, 80, resolver)

	assert.Equal(t, uint(3), calls)
}

func TestLayoutClone(t *testing.T) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("golang"),
		)

	lne.words[0].measured = true

	clone := lne.clone()

	clone.words[0].measured = false

	assert.True(t, lne.words[0].measured)
	assert.False(t, clone.words[0].measured)

	assert.Same(t, lne.frags[0].Base, clone.frags[0].Base)
}

func BenchmarkLayoutMeasure_Cached(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString("hello world"),
	)

	lne.Measure(0, 80)

	b.ReportAllocs()

	for b.Loop() {
		lne.Measure(0, 80)
	}
}

func BenchmarkLayoutMeasure_Recalculate(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString(strings.Repeat("a", 200)),
	)

	cols := winsize.Cols(1)

	b.ReportAllocs()

	for b.Loop() {
		lne.Measure(0, cols)
		cols++
	}
}

func BenchmarkLayoutFindFrags(b *testing.B) {
	lne := emptyLayout().
		PushFrags(
			frag.FromString("a"),
			frag.FromString("b"),
			frag.FromString("c"),
			frag.FromString("d"),
		)

	b.ReportAllocs()

	for b.Loop() {
		_ = lne.FindFrags(0)
	}
}

func BenchmarkLayoutHasAtom(b *testing.B) {
	lne := emptyLayout()

	for range 128 {
		lne.PushFrags(frag.FromString("abc"))
	}

	b.ReportAllocs()

	for b.Loop() {
		lne.HasAtom(0, atom.Break)
	}
}

func BenchmarkLayoutSplitFrag(b *testing.B) {
	frg := frag.FromString(strings.Repeat("a", 200))
	wrd := NewFrag(&frg)

	for b.Loop() {
		splitFragAt(wrd, 40)
	}
}

func BenchmarkSplitLongWord_Fits(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString("hello"),
	)

	b.ReportAllocs()

	for b.Loop() {
		lne.SplitWord(0, 80, 80)
	}
}

func BenchmarkSplitLongWord_SplitMiddle(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString(strings.Repeat("a", 200)),
	)

	b.ReportAllocs()

	for b.Loop() {
		lne.SplitWord(0, 80, 40)
	}
}

func BenchmarkSplitLongWord_SplitFirstRune(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString(strings.Repeat("a", 200)),
	)

	b.ReportAllocs()

	for b.Loop() {
		lne.SplitWord(0, 80, 1)
	}
}

func BenchmarkSplitLongWord_ManyFrags(b *testing.B) {
	frags := make([]frag.Frag, 0, 128)

	for range 128 {
		frags = append(frags, frag.FromString("abcdefghij"))
	}

	lne := emptyLayout().PushFrags(
		frags...,
	)

	b.ReportAllocs()

	for b.Loop() {
		lne.SplitWord(0, 80, 40)
	}
}

func BenchmarkSplitLongWord_WorstCase(b *testing.B) {
	lne := emptyLayout().PushFrags(
		frag.FromString(strings.Repeat("a", 5000)),
	)

	b.ReportAllocs()

	for b.Loop() {
		lne.SplitWord(0, 80, 1)
	}
}
