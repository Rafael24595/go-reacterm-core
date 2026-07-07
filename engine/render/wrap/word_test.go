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
	line := text.LineFromFragments()

	words, frags := splitLineWords(line)

	assert.Size(t, 0, words)
	assert.Size(t, 0, frags)
}

func TestSplitLineWords_EmptyFragmentIgnored(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(""),
		*text.NewFragment("Golang"),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "Golang", wordToString(words[0], frags))
}

func TestSplitLineWords_OnlySpaces(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("   "),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 1, frags)
	assert.Equal(t, "   ", wordToString(words[0], frags))
}

func TestSplitLineWords_StyleChangeRequiresFragmentSplit(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("Zig").AddAtom(atom.Bold),
		*text.NewFragment("lang").AddAtom(atom.Bold),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 2, frags)

	assert.True(t, frags[0].Base.Atom.HasAny(atom.Bold))
	assert.True(t, frags[1].Base.Atom.HasAny(atom.Bold))
}

func TestSplitLineWords_PreservesStylesAcrossFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("ru"),
		*text.NewFragment("st").AddAtom(atom.Select),
		*text.NewFragment("up"),
	)

	words, frags := splitLineWords(line)

	assert.Size(t, 1, words)
	assert.Size(t, 3, frags)

	assert.True(t, frags[0].Base.Atom.HasNone(atom.Select))
	assert.True(t, frags[1].Base.Atom.HasAny(atom.Select))
	assert.True(t, frags[2].Base.Atom.HasNone(atom.Select))
}

func TestSplitLineWords_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(" ").AddAtom(atom.Bold),
		*text.NewFragment(" ").AddAtom(atom.Select),
		*text.NewFragment("c").AddAtom(atom.Bold),
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
	line := text.LineFromFragments(
		*text.NewFragment("c++").AddAtom(atom.Bold),
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

func BenchmarkSplitLineWords_ManyFragments(b *testing.B) {
	frags := make([]text.Fragment, 1000)

	for i := range frags {
		frags[i] = *text.NewFragment("hello ")
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
	line := text.LineFromFragments(
		text.FragmentsFromString(
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
