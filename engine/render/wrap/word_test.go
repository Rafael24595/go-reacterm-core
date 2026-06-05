package wrap

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

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
		*text.NewFragment("Zig").AddAtom(style.AtmBold),
		*text.NewFragment("lang").AddAtom(style.AtmBold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, len(tokens[0].Text))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmBold))
}

func TestSplitLineWords_PreservesStylesAcrossFragments(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("ru"),
		*text.NewFragment("st").AddAtom(style.AtmSelect),
		*text.NewFragment("up"),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasNone(style.AtmSelect))
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmSelect))
	assert.True(t, tokens[0].Text[2].Atom.HasNone(style.AtmSelect))
}

func TestSplitLineWords_MultipleSpaceFragmentsKeepStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment(" ").AddAtom(style.AtmBold),
		*text.NewFragment(" ").AddAtom(style.AtmSelect),
		*text.NewFragment("c").AddAtom(style.AtmBold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 2, len(tokens))

	assert.Equal(t, 2, len(tokens[0].Text))

	assert.Equal(t, " ", tokens[0].Text[0].Text)
	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))

	assert.Equal(t, " ", tokens[0].Text[1].Text)
	assert.True(t, tokens[0].Text[1].Atom.HasAny(style.AtmSelect))

	assert.Equal(t, 1, len(tokens[1].Text))
	assert.Equal(t, "c", tokens[1].Text[0].Text)
	assert.True(t, tokens[1].Text[0].Atom.HasAny(style.AtmBold))
}

func TestSplitLineWords_FinalFlushPreservesStyles(t *testing.T) {
	line := text.LineFromFragments(
		*text.NewFragment("c++").AddAtom(style.AtmBold),
	)

	tokens := splitLineWords(line)

	assert.Equal(t, 1, len(tokens))

	assert.True(t, tokens[0].Text[0].Atom.HasAny(style.AtmBold))
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
