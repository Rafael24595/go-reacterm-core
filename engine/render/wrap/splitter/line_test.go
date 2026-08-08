package splitter

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestSplitLine(t *testing.T) {
	tests := []struct {
		name     string
		line     line.Line
		expected []string
	}{
		{
			name: "single word",
			line: line.FromFrags(
				frag.FromStrings("Golang")...,
			),
			expected: []string{"Golang"},
		},
		{
			name: "word split across frags",
			line: line.FromFrags(
				frag.FromStrings("Z", "ig", "lang")...,
			),
			expected: []string{"Ziglang"},
		},
		{
			name: "two words with space",
			line: line.FromFrags(
				frag.FromStrings("hello cargo")...,
			),
			expected: []string{"hello", " ", "cargo"},
		},
		{
			name: "multiple spaces preserved",
			line: line.FromFrags(
				frag.FromStrings("hello   golangci")...,
			),
			expected: []string{"hello", "   ", "golangci"},
		},
		{
			name: "spaces across frags",
			line: line.FromFrags(
				frag.FromStrings("hello", "  ", "zig")...,
			),
			expected: []string{"hello", "  ", "zig"},
		},
		{
			name: "styled per character",
			line: line.FromFrags(
				frag.FromStrings("r", "u", "s", "t", "c")...,
			),
			expected: []string{"rustc"},
		},
		{
			name: "leading and trailing spaces",
			line: line.FromFrags(
				frag.FromStrings("  Golang  ")...,
			),
			expected: []string{"  ", "Golang", "  "},
		},
		{
			name: "single word across frags",
			line: line.FromFrags(
				frag.FromString("Go"),
				frag.FromString("lang "),
				frag.FromString("Zig"),
				frag.FromString("la"),
				frag.FromString("ng"),
			),
			expected: []string{"Golang", " ", "Ziglang"},
		},
		{
			name: "single long word across frags",
			line: line.FromFrags(
				frag.FromString("supercali"),
				frag.FromString("fragilis"),
				frag.FromString("ticexpia"),
				frag.FromString("lidocious"),
			),
			expected: []string{
				"supercalifragilisticexpialidocious",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words, frags := SplitLine(tt.line)
			got := wordsToStrings(words, frags)

			assert.Size(t, len(tt.expected), got)
			for i := range got {
				assert.Equal(t, tt.expected[i], got[i])
			}
		})
	}
}

func TestSplitLine_EmptyLine(t *testing.T) {
	line := line.FromFrags()

	words, frags := SplitLine(line)

	assert.Empty(t, words)
	assert.Empty(t, frags)
}

func TestSplitLine_EmptyFragIgnored(t *testing.T) {
	line := line.FromFrags(
		frag.FromString(""),
		frag.FromString("Golang"),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "Golang", wordToString(words[0], frags))
}

func TestSplitLine_OnlySpaces(t *testing.T) {
	line := line.FromFrags(
		frag.FromString("   "),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "   ", wordToString(words[0], frags))
}

func TestSplitLine_StyleChangeRequiresFragSplit(t *testing.T) {
	line := line.FromFrags(
		frag.TextAtom("Zig", atom.Bold),
		frag.TextAtom("lang", atom.Bold),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 1, words)
	assert.Size(t, 2, frags)

	assert.True(t, frags[0].Base.Atom().HasAny(atom.Bold))
	assert.True(t, frags[1].Base.Atom().HasAny(atom.Bold))
}

func TestSplitLine_PreservesStylesAcrossFrags(t *testing.T) {
	line := line.FromFrags(
		frag.FromString("ru"),
		frag.TextAtom("st", atom.Select),
		frag.FromString("up"),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 1, words)
	assert.Size(t, 3, frags)

	assert.True(t, frags[0].Base.Atom().HasNone(atom.Select))
	assert.True(t, frags[1].Base.Atom().HasAny(atom.Select))
	assert.True(t, frags[2].Base.Atom().HasNone(atom.Select))
}

func TestSplitLine_MultipleSpaceFragsKeepStyles(t *testing.T) {
	line := line.FromFrags(
		frag.TextAtom(" ", atom.Bold),
		frag.TextAtom(" ", atom.Select),
		frag.TextAtom("c", atom.Bold),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 2, words)
	assert.Size(t, 3, frags)

	word := words[0]
	assert.Size(t, 2, word.End()-word.Start())

	assert.Equal(t, " ", frags[0].Base.Text())
	assert.True(t, frags[0].Base.Atom().HasAny(atom.Bold))

	assert.Equal(t, " ", frags[1].Base.Text())
	assert.True(t, frags[1].Base.Atom().HasAny(atom.Select))

	word = words[1]
	assert.Size(t, 1, word.End()-word.Start())

	assert.Equal(t, "c", frags[2].Base.Text())
	assert.True(t, frags[2].Base.Atom().HasAny(atom.Bold))
}

func TestSplitLine_FinalFlushPreservesStyles(t *testing.T) {
	line := line.FromFrags(
		frag.TextAtom("c++", atom.Bold),
	)

	words, frags := SplitLine(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)

	assert.True(t, frags[0].Base.Atom().HasAny(atom.Bold))
}

func BenchmarkSplitLine(b *testing.B) {
	line := line.FromFrags(
		frag.FromStrings(
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		)...,
	)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = SplitLine(line)
	}
}

func BenchmarkSplitLine_Long(b *testing.B) {
	line := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = SplitLine(line)
	}
}

func BenchmarkSplitLine_ManyFrags(b *testing.B) {
	frags := make([]frag.Frag, 1000)

	for i := range frags {
		frags[i] = frag.FromString("hello ")
	}

	line := line.FromFrags(frags...)

	b.ReportAllocs()

	for b.Loop() {
		SplitLine(line)
	}
}

func BenchmarkSplitWordsCached(b *testing.B) {
	splitter := SplitLineWithCache(
		NewFragCache(),
	)

	l := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		splitter(l)
	}
}

func BenchmarkSplitWords_ColdCache(b *testing.B) {
	cache := NewFragCache()

	splitter := SplitLineWithCache(cache)
	lne := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		cache.Cls()
		_, _ = splitter(lne)
	}
}

func BenchmarkSplitWords_WarmCache(b *testing.B) {
	cache := NewFragCache()

	splitter := SplitLineWithCache(cache)
	lne := benchmarkLine(2000)

	splitter(lne)

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		_, _ = splitter(lne)
	}
}

func BenchmarkSplitLine_ASCII(b *testing.B) {
	line := line.FromString(
		strings.Repeat("hello world ", 300),
	)

	b.ReportAllocs()

	for b.Loop() {
		SplitLine(line)
	}
}

func BenchmarkSplitLine_Unicode(b *testing.B) {
	line := line.FromString(
		strings.Repeat("áéíóú 世界 😀 ", 300),
	)

	b.ReportAllocs()

	for b.Loop() {
		SplitLine(line)
	}
}

func BenchmarkSplitLine_LongWord(b *testing.B) {
	line := line.FromString(
		strings.Repeat("abcdefgh", 1000),
	)

	b.ReportAllocs()

	for b.Loop() {
		SplitLine(line)
	}
}

func BenchmarkSplitLine_ManySpaces(b *testing.B) {
	line := line.FromString(
		strings.Repeat("word     ", 500),
	)

	b.ReportAllocs()

	for b.Loop() {
		SplitLine(line)
	}
}

func BenchmarkSplitLineWith_WarmCache(b *testing.B) {
	cache := NewFragCache()
	splitter := SplitFragWithCache(cache)

	lne := benchmarkLine(2000)

	for _, frg := range lne.Slice() {
		splitter(frg)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		_, _ = SplitLineWith(splitter, lne)
	}
}
