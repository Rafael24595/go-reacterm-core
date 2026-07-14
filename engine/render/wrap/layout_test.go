package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func sourceLayout(source *text.Line) *LayoutLine {
	return &LayoutLine{
		Source: source,
		words:  make([]word, 0),
		frags:  make([]wordFrag, 0),
	}
}

func emptyLayout() *LayoutLine {
	return sourceLayout(&text.Line{})
}

func TestLayoutFindFrags(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("foo"),
			*frag.New("bar"),
		).
		pushFrags(
			*frag.New("baz"),
		)

	frags := layout.findFrags(0)

	assert.Equal(t, 2, len(frags))
	assert.Equal(t, "foo", frags[0].Base.Text)
	assert.Equal(t, "bar", frags[1].Base.Text)

	last := layout.findFrags(1)

	assert.Equal(t, 1, len(last))
	assert.Equal(t, "baz", last[0].Base.Text)
}

func TestLayoutPushFrags(t *testing.T) {
	layout := emptyLayout()

	layout.
		pushFrags(*frag.New("a")).
		pushFrags(
			*frag.New("b"),
			*frag.New("c"),
		)

	assert.Equal(t, 2, len(layout.words))
	assert.Equal(t, 3, len(layout.frags))

	assert.Equal(t, uint32(0), layout.words[0].start)
	assert.Equal(t, uint32(1), layout.words[0].end)

	assert.Equal(t, uint32(1), layout.words[1].start)
	assert.Equal(t, uint32(3), layout.words[1].end)
}

func TestLayoutHasAtom(t *testing.T) {
	frg := frag.New("foo")
	frg.Atom = atom.Wrap

	layout := emptyLayout().
		pushFrags(*frg)

	assert.True(t, layout.hasAtom(0, atom.Wrap))
	assert.False(t, layout.hasAtom(0, atom.Focus))
}

func TestSplitWord(t *testing.T) {
	tests := []struct {
		name            string
		layout          *LayoutLine
		cols            winsize.Cols
		remaining       winsize.Cols
		expectedCurrent string
		expectedRest    string
	}{
		{
			name: "word fits completely",
			layout: emptyLayout().pushFrags(
				*frag.New("golang"),
			),
			cols:            20,
			remaining:       20,
			expectedCurrent: "golang",
			expectedRest:    "",
		},
		{
			name: "split single frag word",
			layout: emptyLayout().pushFrags(
				*frag.New("ziglang"),
			),
			cols:            4,
			remaining:       4,
			expectedCurrent: "zigl",
			expectedRest:    "ang",
		},
		{
			name: "split fragmented word",
			layout: emptyLayout().pushFrags(
				*frag.New("go"),
				*frag.New("la"),
				*frag.New("ng"),
			),
			cols:            2,
			remaining:       4,
			expectedCurrent: "gola",
			expectedRest:    "ng",
		},
		{
			name: "zero remaining",
			layout: emptyLayout().pushFrags(
				*frag.New("rust"),
			),
			cols:            5,
			remaining:       0,
			expectedCurrent: "",
			expectedRest:    "rust",
		},
		{
			name: "split inside second frag",
			layout: emptyLayout().pushFrags(
				*frag.New("cl"),
				*frag.New("oju"),
				*frag.New("re"),
			),
			cols:            3,
			remaining:       3,
			expectedCurrent: "clo",
			expectedRest:    "jure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := 0

			ok := tt.layout.splitWord(
				idx, tt.cols, tt.remaining,
			)

			if tt.expectedCurrent != "" {
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCurrent, fragsToString(tt.layout.findFrags(idx)))
			}
		})
	}
}

func TestSplitWord_FitsWithoutMutatingLayout(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("golang"),
		)

	ok := layout.splitWord(0, 80, 80)

	assert.True(t, ok)

	assert.Size(t, 1, layout.words)
	assert.Size(t, 1, layout.frags)

	assert.Equal(t, "golang", fragsToString(layout.findFrags(0)))

	assert.Equal(t, uint32(0), layout.words[0].start)
	assert.Equal(t, uint32(1), layout.words[0].end)
}

func TestSplitWord_SplitLastFrag(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("go"),
			*frag.New("la"),
			*frag.New("ng"),
		)

	ok := layout.splitWord(0, 80, 5)

	assert.True(t, ok)

	assert.Size(t, 2, layout.words)
	assert.Size(t, 4, layout.frags)

	assert.Equal(t, "golan", fragsToString(layout.findFrags(0)))
	assert.Equal(t, "n", layout.frags[2].Base.Text)
	assert.Equal(t, "g", fragsToString(layout.findFrags(1)))
}

func TestLayoutSplitFrag(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("abcdef"),
		)

	layout.splitFrag(0, 0, 3)

	assert.Equal(t, 2, len(layout.frags))
	assert.Equal(t, 2, len(layout.words))

	assert.Equal(t, "abc", layout.frags[0].Base.Text)
	assert.Equal(t, "def", layout.frags[1].Base.Text)

	assert.Equal(t, uint32(0), layout.words[0].start)
	assert.Equal(t, uint32(1), layout.words[0].end)

	assert.Equal(t, uint32(1), layout.words[1].start)
	assert.Equal(t, uint32(2), layout.words[1].end)
}

func TestLayoutSplitFrag_InvalidatesMeasureCache(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("abcdef"),
		)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return winsize.Cols(len(frags))
	}

	layout.measureWith(0, 80, resolver)
	assert.Equal(t, uint(1), calls)

	layout.splitFrag(0, 0, 3)

	layout.measureWith(0, 80, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestLayoutSplitFrag_ShiftsFollowingWords(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("abcdef"),
		).
		pushFrags(
			*frag.New("foo"),
		).
		pushFrags(
			*frag.New("bar"),
		)

	layout.splitFrag(0, 0, 3)

	assert.Size(t, 4, layout.words)
	assert.Size(t, 4, layout.frags)

	assert.Equal(t, "abc", fragsToString(layout.findFrags(0)))
	assert.Equal(t, "def", fragsToString(layout.findFrags(1)))
	assert.Equal(t, "foo", fragsToString(layout.findFrags(2)))
	assert.Equal(t, "bar", fragsToString(layout.findFrags(3)))

	assert.Equal(t, 0, layout.words[0].start)
	assert.Equal(t, 1, layout.words[0].end)

	assert.Equal(t, 1, layout.words[1].start)
	assert.Equal(t, 2, layout.words[1].end)

	assert.Equal(t, 2, layout.words[2].start)
	assert.Equal(t, 3, layout.words[2].end)

	assert.Equal(t, 3, layout.words[3].start)
	assert.Equal(t, 4, layout.words[3].end)
}

func TestLayoutSplitFrag_NoSplit(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("abc"),
		)

	layout.splitFrag(0, 0, 3)

	assert.Size(t, 1, layout.words)
	assert.Size(t, 1, layout.frags)

	assert.Equal(t, "abc", fragsToString(layout.findFrags(0)))
}

func TestLayoutWordMeasure_CacheSameCols(t *testing.T) {
	layout := emptyLayout().pushFrags(
		*frag.New("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	first := layout.measureWith(0, 80, resolver)
	second := layout.measureWith(0, 80, resolver)

	assert.Equal(t, first, second)
	assert.Equal(t, winsize.Cols(80), layout.words[0].cols)

	assert.Equal(t, 1, calls)
}

func TestLayoutWordMeasure_RecalculateOnColsChange(t *testing.T) {
	layout := emptyLayout().pushFrags(
		*frag.New("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	_ = layout.measureWith(0, 80, resolver)
	m40 := layout.measureWith(0, 40, resolver)

	assert.Equal(t, winsize.Cols(40), layout.words[0].cols)
	assert.Equal(t, m40, layout.words[0].measure)

	assert.Equal(t, 2, calls)
}

func TestLayoutWordMeasure_CacheAfterColsChange(t *testing.T) {
	layout := emptyLayout().pushFrags(
		*frag.New("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	layout.measureWith(0, 80, resolver)
	layout.measureWith(0, 40, resolver)
	layout.measureWith(0, 40, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestLayoutWordMeasure_RecalculateWhenReturningToPreviousCols(t *testing.T) {
	layout := emptyLayout().pushFrags(
		*frag.New("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...wordFrag) winsize.Cols {
		calls++
		return 42
	}

	layout.measureWith(0, 80, resolver)
	layout.measureWith(0, 40, resolver)
	layout.measureWith(0, 80, resolver)

	assert.Equal(t, uint(3), calls)
}

func TestLayoutClone(t *testing.T) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("golang"),
		)

	layout.words[0].measured = true

	clone := layout.clone()

	clone.words[0].measured = false

	assert.True(t, layout.words[0].measured)
	assert.False(t, clone.words[0].measured)

	clone.frags[0].Base.Text = "rust"

	assert.Equal(t, "rust", layout.frags[0].Base.Text)
	assert.Equal(t, "rust", clone.frags[0].Base.Text)
}

func BenchmarkLayoutMeasure_Cached(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New("hello world"),
	)

	layout.measure(0, 80)

	b.ReportAllocs()

	for b.Loop() {
		layout.measure(0, 80)
	}
}

func BenchmarkLayoutMeasure_Recalculate(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New(strings.Repeat("a", 200)),
	)

	cols := winsize.Cols(1)

	b.ReportAllocs()

	for b.Loop() {
		layout.measure(0, cols)
		cols++
	}
}

func BenchmarkLayoutFindFrags(b *testing.B) {
	layout := emptyLayout().
		pushFrags(
			*frag.New("a"),
			*frag.New("b"),
			*frag.New("c"),
			*frag.New("d"),
		)

	b.ReportAllocs()

	for b.Loop() {
		_ = layout.findFrags(0)
	}
}

func BenchmarkLayoutHasAtom(b *testing.B) {
	layout := emptyLayout()

	for range 128 {
		layout.pushFrags(*frag.New("abc"))
	}

	b.ReportAllocs()

	for b.Loop() {
		layout.hasAtom(0, atom.Break)
	}
}

func BenchmarkLayoutSplitFrag(b *testing.B) {
	frag := newWordFrag(
		frag.New(strings.Repeat("a", 200)),
	)

	for b.Loop() {
		splitFragAt(frag, 40)
	}
}

func BenchmarkSplitLongWord_Fits(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New("hello"),
	)

	b.ReportAllocs()

	for b.Loop() {
		layout.splitWord(0, 80, 80)
	}
}

func BenchmarkSplitLongWord_SplitMiddle(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New(strings.Repeat("a", 200)),
	)

	b.ReportAllocs()

	for b.Loop() {
		layout.splitWord(0, 80, 40)
	}
}

func BenchmarkSplitLongWord_SplitFirstRune(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New(strings.Repeat("a", 200)),
	)

	b.ReportAllocs()

	for b.Loop() {
		layout.splitWord(0, 80, 1)
	}
}

func BenchmarkSplitLongWord_ManyFrags(b *testing.B) {
	frags := make([]frag.Frag, 0, 128)

	for range 128 {
		frags = append(frags, *frag.New("abcdefghij"))
	}

	layout := emptyLayout().pushFrags(
		frags...,
	)

	b.ReportAllocs()

	for b.Loop() {
		layout.splitWord(0, 80, 40)
	}
}

func BenchmarkSplitLongWord_WorstCase(b *testing.B) {
	layout := emptyLayout().pushFrags(
		*frag.New(strings.Repeat("a", 5000)),
	)

	b.ReportAllocs()

	for b.Loop() {
		layout.splitWord(0, 80, 1)
	}
}
