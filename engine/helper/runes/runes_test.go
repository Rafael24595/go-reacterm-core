package runes

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestJoinReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "basic",
			in:   []string{"a", "b", "c"},
			out:  "cba",
		},
		{
			name: "words",
			in:   []string{"hello", " ", "golang"},
			out:  "golang hello",
		},
		{
			name: "unicode",
			in:   []string{"🙂", "🚀", "go"},
			out:  "go🚀🙂",
		},
		{
			name: "empty",
			in:   []string{},
			out:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, JoinReverse(tt.in))
		})
	}
}

func TestRuneIndexToByteIndex(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		runeIndex offset.Offset
		expected  offset.Offset
		ok        bool
	}{
		{
			name:      "ascii simple",
			text:      "hello",
			runeIndex: 1,
			expected:  1,
			ok:        true,
		},
		{
			name:      "unicode multi-byte",
			text:      "a🙂b",
			runeIndex: 1,
			expected:  1,
			ok:        true,
		},
		{
			name:      "unicode end",
			text:      "a🙂b",
			runeIndex: 3,
			expected:  offset.Offset(len("a🙂b")),
			ok:        true,
		},
		{
			name:      "zero index",
			text:      "abc",
			runeIndex: 0,
			expected:  0,
			ok:        true,
		},
		{
			name:      "out of bounds",
			text:      "abc",
			runeIndex: 5,
			expected:  0,
			ok:        false,
		},
		{
			name:      "exact end boundary",
			text:      "abc",
			runeIndex: 3,
			expected:  3,
			ok:        true,
		},
		{
			name:      "empty string",
			text:      "",
			runeIndex: 0,
			expected:  0,
			ok:        true,
		},
		{
			name:      "multi rune unicode",
			text:      "🙂🙂🙂",
			runeIndex: 2,
			expected:  offset.Offset(len("🙂🙂")),
			ok:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, ok := RuneIndexToByteIndex(tt.text, tt.runeIndex)

			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, idx)
		})
	}
}

func TestSanitizeRunes(t *testing.T) {
	tests := []struct {
		name     string
		input    []rune
		want     []rune
		wantSame bool
	}{
		{
			name:     "Without null values",
			input:    []rune{'a', 'b', 'c'},
			want:     []rune{'a', 'b', 'c'},
			wantSame: true,
		},
		{
			name:     "Empty buffer",
			input:    []rune{},
			want:     []rune{},
			wantSame: true,
		},
		{
			name:     "Nil buffer",
			input:    nil,
			want:     nil,
			wantSame: true,
		},
		{
			name:     "Null at start",
			input:    []rune{0, 'a', 'b'},
			want:     []rune{'a', 'b'},
			wantSame: false,
		},
		{
			name:     "Null at end",
			input:    []rune{'a', 'b', 0},
			want:     []rune{'a', 'b'},
			wantSame: false,
		},
		{
			name:     "Null intercalated",
			input:    []rune{'a', 0, 'b', 0, 'c'},
			want:     []rune{'a', 'b', 'c'},
			wantSame: false,
		},
		{
			name:     "Only",
			input:    []rune{0, 0, 0},
			want:     []rune{},
			wantSame: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeRunes(tt.input)

			assert.DeepEqual(t, tt.want, got)

			if tt.wantSame && len(tt.input) > 0 {
				assert.Equal(t, &tt.input[0], &got[0])
			}
		})
	}
}

func TestMeasure(t *testing.T) {
	tests := []struct {
		name string
		text string
		want winsize.Cols
	}{
		{"ascii", "hello", 5},
		{"unicode", "🙂🙂", 2},
		{"mixed", "a🙂b", 3},
		{"empty", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Measure(tt.text))
		})
	}
}
