package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func wordsToString(words ...word) string {
	var sb strings.Builder

	for _, w := range words {
		for _, f := range w.Text {
			sb.WriteString(f.Base.Text)
		}
	}

	return sb.String()
}

func tokenString(token word) string {
	var b strings.Builder
	for _, f := range token.Text {
		b.WriteString(f.Text)
	}
	return b.String()
}

func tokenStrings(tokens []word) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = tokenString(t)
	}
	return out
}

func TestSplitLineWords(t *testing.T) {
	tests := []struct {
		name     string
		line     *text.Line
		expected []string
	}{
		{
			name: "single word",
			line: text.LineFromFragments(
				text.FragmentsFromString("Golang")...,
			),
			expected: []string{"Golang"},
		},
		{
			name: "word split across fragments",
			line: text.LineFromFragments(
				text.FragmentsFromString("Z", "ig", "lang")...,
			),
			expected: []string{"Ziglang"},
		},
		{
			name: "two words with space",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello cargo")...,
			),
			expected: []string{"hello", " ", "cargo"},
		},
		{
			name: "multiple spaces preserved",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello   golangci")...,
			),
			expected: []string{"hello", "   ", "golangci"},
		},
		{
			name: "spaces across fragments",
			line: text.LineFromFragments(
				text.FragmentsFromString("hello", "  ", "zig")...,
			),
			expected: []string{"hello", "  ", "zig"},
		},
		{
			name: "styled per character",
			line: text.LineFromFragments(
				text.FragmentsFromString("r", "u", "s", "t", "c")...,
			),
			expected: []string{"rustc"},
		},
		{
			name: "leading and trailing spaces",
			line: text.LineFromFragments(
				text.FragmentsFromString("  Golang  ")...,
			),
			expected: []string{"  ", "Golang", "  "},
		},
		{
			name: "single word across fragments",
			line: text.LineFromFragments(
				*text.NewFragment("Go"),
				*text.NewFragment("lang "),
				*text.NewFragment("Zig"),
				*text.NewFragment("la"),
				*text.NewFragment("ng"),
			),
			expected: []string{"Golang", " ", "Ziglang"},
		},
		{
			name: "single long word across fragments",
			line: text.LineFromFragments(
				*text.NewFragment("supercali"),
				*text.NewFragment("fragilis"),
				*text.NewFragment("ticexpia"),
				*text.NewFragment("lidocious"),
			),
			expected: []string{
				"supercalifragilisticexpialidocious",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := splitLineWords(tt.line)
			got := tokenStrings(tokens)

			assert.Size(t, len(tt.expected), got)
			for i := range got {
				assert.Equal(t, tt.expected[i], got[i])
			}
		})
	}
}

func TestSplitLineWords_EmptyLine(t *testing.T) {
	line := text.LineFromFragments()

	tokens := splitLineWords(line)
	assert.Equal(t, 0, len(tokens))
}

func TestSplitLineWords_EmptyFragmentIgnored(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(""),
		*text.NewFragment("Golang"),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "Golang", tokenString(tokens[0]))
}

func TestSplitLineWords_OnlySpaces(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("   "),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, "   ", tokenString(tokens[0]))
}

func TestSplitLineWords_StyleChangeRequiresFragmentSplit(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("Zig").AddAtom(atom.Bold),
		*text.NewFragment("lang").AddAtom(atom.Bold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, len(tokens[0].Text))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(atom.Bold))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(atom.Bold))
}

func TestSplitLineWords_PreservesStylesAcrossFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("ru"),
		*text.NewFragment("st").AddAtom(atom.Select),
		*text.NewFragment("up"),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasNone(atom.Select))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(atom.Select))
	assert.True(t, tokens[0].Text[2].Atom.HasNone(atom.Select))
}

func TestSplitLineWords_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(" ").AddAtom(atom.Bold),
		*text.NewFragment(" ").AddAtom(atom.Select),
		*text.NewFragment("c").AddAtom(atom.Bold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 2, len(tokens))

	assert.Equal(t, 2, len(tokens[0].Text))

	assert.Equal(t, " ", tokens[0].Text[0].Text)
	assert.True(t, tokens[0].Text[0].Atom.HasAny(atom.Bold))

	assert.Equal(t, " ", tokens[0].Text[1].Text)
	assert.True(t, tokens[0].Text[1].Atom.HasAny(atom.Select))

	assert.Equal(t, 1, len(tokens[1].Text))
	assert.Equal(t, "c", tokens[1].Text[0].Text)
	assert.True(t, tokens[1].Text[0].Atom.HasAny(atom.Bold))
}

func TestSplitLineWords_FinalFlushPreservesStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("c++").AddAtom(atom.Bold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(atom.Bold))
}

func TestSplitLongWord(t *testing.T) {
	tests := []struct {
		name            string
		word            word
		cols            winsize.Cols
		remaining       winsize.Cols
		expectedCurrent string
		expectedRest    string
	}{
		{
			name: "word fits completely",
			word: *newWord(
				*text.NewFragment("golang"),
			),
			cols:            20,
			remaining:       20,
			expectedCurrent: "golang",
			expectedRest:    "",
		},
		{
			name: "split single fragment word",
			word: *newWord(
				*text.NewFragment("ziglang"),
			),
			cols:            4,
			remaining:       4,
			expectedCurrent: "zigl",
			expectedRest:    "ang",
		},
		{
			name: "split fragmented word",
			word: *newWord(
				*text.NewFragment("go"),
				*text.NewFragment("la"),
				*text.NewFragment("ng"),
			),
			cols:            2,
			remaining:       4,
			expectedCurrent: "gola",
			expectedRest:    "ng",
		},
		{
			name: "zero remaining",
			word: *newWord(
				*text.NewFragment("rust"),
			),
			cols:            5,
			remaining:       0,
			expectedCurrent: "",
			expectedRest:    "rust",
		},
		{
			name: "split inside second fragment",
			word: *newWord(
				*text.NewFragment("cl"),
				*text.NewFragment("oju"),
				*text.NewFragment("re"),
			),
			cols:            3,
			remaining:       3,
			expectedCurrent: "clo",
			expectedRest:    "jure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current, rest := splitLongWord(
				tt.word,
				tt.cols,
				tt.remaining,
			)

			if tt.expectedCurrent != "" {
				assert.NotNil(t, current)
				assert.Equal(t, tt.expectedCurrent, tokenString(*current))
			}

			if tt.expectedRest != "" {
				assert.NotNil(t, rest)
				assert.Equal(t, tt.expectedRest, tokenString(*rest))
			}
		})
	}
}

func TestWordMeasure_CacheSameCols(t *testing.T) {
	w := newWord(
		*text.NewFragment("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...text.Fragment) winsize.Cols {
		calls++
		return 42
	}

	first := w.measureWith(80, resolver)
	second := w.measureWith(80, resolver)

	assert.Equal(t, first, second)
	assert.Equal(t, winsize.Cols(80), w.cols)

	assert.Equal(t, 1, calls)
}

func TestWordMeasure_RecalculateOnColsChange(t *testing.T) {
	w := newWord(
		*text.NewFragment("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...text.Fragment) winsize.Cols {
		calls++
		return 42
	}

	_ = w.measureWith(80, resolver)
	m40 := w.measureWith(40, resolver)

	assert.Equal(t, winsize.Cols(40), w.cols)
	assert.Equal(t, m40, w.measure)

	assert.Equal(t, 2, calls)
}

func TestWordMeasure_CacheAfterColsChange(t *testing.T) {
	w := newWord(
		*text.NewFragment("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...text.Fragment) winsize.Cols {
		calls++
		return 42
	}

	w.measureWith(80, resolver)
	w.measureWith(40, resolver)
	w.measureWith(40, resolver)

	assert.Equal(t, uint(2), calls)
}

func TestWordMeasure_RecalculateWhenReturningToPreviousCols(t *testing.T) {
	w := newWord(
		*text.NewFragment("golang"),
	)

	calls := uint(0)

	resolver := func(cols winsize.Cols, frags ...text.Fragment) winsize.Cols {
		calls++
		return 42
	}

	w.measureWith(80, resolver)
	w.measureWith(40, resolver)
	w.measureWith(80, resolver)

	assert.Equal(t, uint(3), calls)
}

func BenchmarkSplitLineWords(b *testing.B) {
	line := text.LineFromFragments(
		text.FragmentsFromString(
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. "+
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		)...,
	)

	b.ReportAllocs()
	

	for b.Loop() {
		_ = splitLineWords(line)
	}
}

func BenchmarkSplitLineWords_Long(b *testing.B) {
    line := benchmarkLine(2000)

    b.ReportAllocs()
    

    for b.Loop() {
        _ = splitLineWords(&line)
    }
}

func BenchmarkSplitLongWord_Fits(b *testing.B) {
    w := wordFromFrags(
        toWordFrag(
            *text.NewFragment("hello"),
        )...,
    )

    b.ReportAllocs()
    

    for b.Loop() {
        splitLongWord(*w, 80, 80)
    }
}

func BenchmarkSplitLongWord_SplitMiddle(b *testing.B) {
    w := wordFromFrags(
        toWordFrag(
            *text.NewFragment(strings.Repeat("a", 200)),
        )...,
    )

    b.ReportAllocs()
    

    for b.Loop() {
        splitLongWord(*w, 80, 40)
    }
}

func BenchmarkSplitLongWord_SplitFirstRune(b *testing.B) {
    w := wordFromFrags(
        toWordFrag(
            *text.NewFragment(strings.Repeat("a", 200)),
        )...,
    )

    b.ReportAllocs()
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        splitLongWord(*w, 80, 1)
    }
}

func BenchmarkSplitLongWord_ManyFragments(b *testing.B) {
    frags := make([]text.Fragment, 0, 128)

    for range 128 {
        frags = append(frags, *text.NewFragment("abcdefghij"))
    }

    w := wordFromFrags(
        toWordFrag(frags...)...,
    )

    b.ReportAllocs()
    

    for b.Loop() {
        splitLongWord(*w, 80, 40)
    }
}

func BenchmarkSplitLongWord_WorstCase(b *testing.B) {
    w := wordFromFrags(
        toWordFrag(
            *text.NewFragment(strings.Repeat("a", 5000)),
        )...,
    )

    b.ReportAllocs()
    

    for b.Loop() {
        splitLongWord(*w, 80, 1)
    }
}
