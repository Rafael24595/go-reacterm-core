package runes

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestInsert(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []rune
		insert   []rune
		pos      offset.Offset
		expected string
	}{
		{
			name:     "Insert at start (pos 0)",
			buffer:   []rune{'w', 'o', 'r', 'l', 'd'},
			insert:   []rune{'h', 'e', 'l', 'l', 'o', ' '},
			pos:      0,
			expected: "hello world",
		},
		{
			name:     "Insert in middle",
			buffer:   []rune{'a', 'b', 'e'},
			insert:   []rune{'c', 'd'},
			pos:      2,
			expected: "abcde",
		},
		{
			name:     "Insert at exact end",
			buffer:   []rune{'g', 'o'},
			insert:   []rune{'l', 'a', 'n', 'g'},
			pos:      2,
			expected: "golang",
		},
		{
			name:     "Insert into empty buffer",
			buffer:   []rune{},
			insert:   []rune{'a', 'b', 'c'},
			pos:      0,
			expected: "abc",
		},
		{
			name:     "Insert empty slice (no-op)",
			buffer:   []rune{'k', 'e', 'e', 'p'},
			insert:   []rune{},
			pos:      2,
			expected: "keep",
		},
		{
			name:     "Insert Unicode characters",
			buffer:   []rune{'a', 'c'},
			insert:   []rune{'🚀'},
			pos:      1,
			expected: "a🚀c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Insert(tt.buffer, tt.insert, tt.pos)
			assert.Equal(t, tt.expected, string(got))
		})
	}
}

func TestInsert_Immutability(t *testing.T) {
	original := []rune{'a', 'b', 'd'}
	originalCopy := append([]rune(nil), original...)

	insert := []rune{'c'}

	result := Insert(original, insert, 2)

	result[0] = 'Z'

	assert.Equal(t, "abd", string(originalCopy))
	assert.Equal(t, "abd", string(original))
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []rune
		insert   []rune
		start    offset.Offset
		end      offset.Offset
		expected string
	}{
		{
			name:     "Replace sub-slice in middle (same length)",
			buffer:   []rune{'h', 'e', 'x', 'x', 'o'},
			insert:   []rune{'l', 'l'},
			start:    2,
			end:      4,
			expected: "hello",
		},
		{
			name:     "Replace range with smaller slice (shrinking buffer)",
			buffer:   []rune{'a', 'b', 'c', 'd', 'e', 'f'},
			insert:   []rune{'X'},
			start:    1,
			end:      5,
			expected: "aXf",
		},
		{
			name:     "Replace range with larger slice (expanding buffer)",
			buffer:   []rune{'a', 'b', 'f'},
			insert:   []rune{'c', 'd', 'e'},
			start:    2,
			end:      2, // Zero range, acts like AppendAt
			expected: "abcdef",
		},
		{
			name:     "Replace entire buffer",
			buffer:   []rune{'o', 'l', 'd'},
			insert:   []rune{'n', 'e', 'w'},
			start:    0,
			end:      3,
			expected: "new",
		},
		{
			name:     "Delete range by passing empty insert",
			buffer:   []rune{'a', 'b', 'C', 'D', 'e'},
			insert:   []rune{},
			start:    2,
			end:      4,
			expected: "abe",
		},
		{
			name:     "Replace prefix range",
			buffer:   []rune{'b', 'a', 'd', 'c', 'a', 't'},
			insert:   []rune{'g', 'o', 'o', 'd'},
			start:    0,
			end:      3,
			expected: "goodcat",
		},
		{
			name:     "Replace suffix range up to buffer end",
			buffer:   []rune{'c', 'o', 'd', 'e', 'X', 'Y'},
			insert:   []rune{'r'},
			start:    4,
			end:      6,
			expected: "coder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Replace(tt.buffer, tt.insert, tt.start, tt.end)
			assert.Equal(t, tt.expected, string(got))
		})
	}
}

func TestReplace_Immutability(t *testing.T) {
	original := []rune{'h', 'e', 'X', 'X', 'o'}
	originalCopy := append([]rune(nil), original...)

	insert := []rune{'l', 'l'}

	result := Replace(original, insert, 2, 4)

	result[2] = 'Z'

	assert.Equal(t, "heXXo", string(originalCopy))
	assert.Equal(t, "heXXo", string(original))
}

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

func TestMeasureGenerics(t *testing.T) {
	text := "a🙂b"
	runes := []rune{'a', '🙂', 'b'}

	assert.Equal(t, 3, Measure[int](text))
	assert.Equal(t, winsize.Cols(3), Measure[winsize.Cols](text))
	assert.Equal(t, offset.Offset(3), Measure[offset.Offset](text))

	assert.Equal(t, 3, MeasureRunes[int](runes))
	assert.Equal(t, winsize.Cols(3), MeasureRunes[winsize.Cols](runes))
	assert.Equal(t, offset.Offset(3), MeasureRunes[offset.Offset](runes))
}

func TestMeasureWrappers(t *testing.T) {
	text := "hello"
	runes := []rune{'h', 'e', 'l', 'l', 'o'}

	assert.Equal(t, winsize.Cols(5), MeasureCols(text))
	assert.Equal(t, winsize.Cols(5), MeasureColsRunes(runes))

	assert.Equal(t, offset.Offset(5), MeasureOffset(text))
	assert.Equal(t, offset.Offset(5), MeasureOffsetRunes(runes))
}
