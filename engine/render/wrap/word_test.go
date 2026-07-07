package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func wordsToString(words []word, frags []wordFrag) string {
	var sb strings.Builder
	for _, words := range words {
		sb.WriteString(
			wordToString(words, frags),
		)
	}
	return sb.String()
}

func wordToString(word word, frags []wordFrag) string {
	return fragsToString(frags[word.start:word.end])
}

func wordsToStrings(tokens []word, frags []wordFrag) []string {
	out := make([]string, len(tokens))
	for i, word := range tokens {
		out[i] = wordToString(word, frags)
	}
	return out
}

func fragsToString(frags []wordFrag) string {
	var b strings.Builder
	for _, f := range frags {
		b.WriteString(f.Base.Text)
	}
	return b.String()
}

func TestSplitLineWords(t *testing.T) {
	tests := []struct {
		name     string
		line     *text.Line
		expected []string
	}{
		{
			name: "single word",
			line: text.LineFromFrags(
				text.FragsFromString("Golang")...,
			),
			expected: []string{"Golang"},
		},
		{
			name: "word split across frags",
			line: text.LineFromFrags(
				text.FragsFromString("Z", "ig", "lang")...,
			),
			expected: []string{"Ziglang"},
		},
		{
			name: "two words with space",
			line: text.LineFromFrags(
				text.FragsFromString("hello cargo")...,
			),
			expected: []string{"hello", " ", "cargo"},
		},
		{
			name: "multiple spaces preserved",
			line: text.LineFromFrags(
				text.FragsFromString("hello   golangci")...,
			),
			expected: []string{"hello", "   ", "golangci"},
		},
		{
			name: "spaces across frags",
			line: text.LineFromFrags(
				text.FragsFromString("hello", "  ", "zig")...,
			),
			expected: []string{"hello", "  ", "zig"},
		},
		{
			name: "styled per character",
			line: text.LineFromFrags(
				text.FragsFromString("r", "u", "s", "t", "c")...,
			),
			expected: []string{"rustc"},
		},
		{
			name: "leading and trailing spaces",
			line: text.LineFromFrags(
				text.FragsFromString("  Golang  ")...,
			),
			expected: []string{"  ", "Golang", "  "},
		},
		{
			name: "single word across frags",
			line: text.LineFromFrags(
				*text.NewFrag("Go"),
				*text.NewFrag("lang "),
				*text.NewFrag("Zig"),
				*text.NewFrag("la"),
				*text.NewFrag("ng"),
			),
			expected: []string{"Golang", " ", "Ziglang"},
		},
		{
			name: "single long word across frags",
			line: text.LineFromFrags(
				*text.NewFrag("supercali"),
				*text.NewFrag("fragilis"),
				*text.NewFrag("ticexpia"),
				*text.NewFrag("lidocious"),
			),
			expected: []string{
				"supercalifragilisticexpialidocious",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			words, frags := splitLineWords(tt.line)
			got := wordsToStrings(words, frags)

			assert.Size(t, len(tt.expected), got)
			for i := range got {
				assert.Equal(t, tt.expected[i], got[i])
			}
		})
	}
}

func TestSplitLineWords_EmptyLine(t *testing.T) {
	line := text.LineFromFrags()

	words, frags := splitLineWords(line)

	assert.Size(t, 0, words)
	assert.Size(t, 0, frags)
}

func TestSplitLineWords_EmptyFragIgnored(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag(""),
		*text.NewFrag("Golang"),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "Golang", wordToString(words[0], frags))
}

func TestSplitLineWords_OnlySpaces(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("   "),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "   ", wordToString(words[0], frags))
}

func TestSplitLineWords_StyleChangeRequiresFragSplit(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("Zig").AddAtom(atom.Bold),
		*text.NewFrag("lang").AddAtom(atom.Bold),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 2, frags)

	assert.True(t, frags[0].Base.Atom.HasAny(atom.Bold))
	assert.True(t, frags[1].Base.Atom.HasAny(atom.Bold))
}

func TestSplitLineWords_PreservesStylesAcrossFrags(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("ru"),
		*text.NewFrag("st").AddAtom(atom.Select),
		*text.NewFrag("up"),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 3, frags)

	assert.True(t, frags[0].Base.Atom.HasNone(atom.Select))
	assert.True(t, frags[1].Base.Atom.HasAny(atom.Select))
	assert.True(t, frags[2].Base.Atom.HasNone(atom.Select))
}

func TestSplitLineWords_MultipleSpaceFragsKeepStyles(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag(" ").AddAtom(atom.Bold),
		*text.NewFrag(" ").AddAtom(atom.Select),
		*text.NewFrag("c").AddAtom(atom.Bold),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 2, words)
	assert.Size(t, 3, frags)

	word := words[0]
	assert.Size(t, 2, word.end-word.start)

	assert.Equal(t, " ", frags[0].Base.Text)
	assert.True(t, frags[0].Base.Atom.HasAny(atom.Bold))

	assert.Equal(t, " ", frags[1].Base.Text)
	assert.True(t, frags[1].Base.Atom.HasAny(atom.Select))

	word = words[1]
	assert.Size(t, 1, word.end-word.start)

	assert.Equal(t, "c", frags[2].Base.Text)
	assert.True(t, frags[2].Base.Atom.HasAny(atom.Bold))
}

func TestSplitLineWords_FinalFlushPreservesStyles(t *testing.T) {
	line := text.LineFromFrags(
		*text.NewFrag("c++").AddAtom(atom.Bold),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)

	assert.True(t, frags[0].Base.Atom.HasAny(atom.Bold))
}

func BenchmarkSplitLineFeeds_NoLF(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("Hello World ", 100),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineFeeds(&line, false)
	}
}

func BenchmarkSplitLineFeeds_SomeLF(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("Hello\nWorld\n", 100),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineFeeds(&line, false)
	}
}

func BenchmarkSplitLineFeeds_ManyLF(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("\n", 1000),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineFeeds(&line, false)
	}
}

func BenchmarkSplitLineWords_ASCII(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("hello world ", 300),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineWords(&line)
	}
}

func BenchmarkSplitLineWords_Unicode(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("áéíóú 世界 😀 ", 300),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineWords(&line)
	}
}

func BenchmarkSplitLineWords_LongWord(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("abcdefgh", 1000),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineWords(&line)
	}
}

func BenchmarkSplitLineWords_ManySpaces(b *testing.B) {
	line := *text.NewLine(
		strings.Repeat("word     ", 500),
	)

	b.ReportAllocs()

	for b.Loop() {
		splitLineWords(&line)
	}
}

func BenchmarkSplitLineWords_ManyFrags(b *testing.B) {
	frags := make([]text.Frag, 1000)

	for i := range frags {
		frags[i] = *text.NewFrag("hello ")
	}

	line := text.Line{
		Text: frags,
	}

	b.ReportAllocs()

	for b.Loop() {
		splitLineWords(&line)
	}
}

func BenchmarkSplitLineWords(b *testing.B) {
	line := text.LineFromFrags(
		text.FragsFromString(
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
				"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		)...,
	)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = splitLineWords(line)
	}
}

func BenchmarkSplitLineWords_Long(b *testing.B) {
	line := benchmarkLine(2000)

	b.ReportAllocs()

	for b.Loop() {
		_, _ = splitLineWords(&line)
	}
}
