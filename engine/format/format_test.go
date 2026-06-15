package format

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestAlignCenter(t *testing.T) {
	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		filler   string
		expected string
	}{
		{"already fits", 3, TextFromString("abc"), "-", "abc"},
		{"smaller width", 2, TextFromString("abc"), "-", "abc"},
		{"even padding", 7, TextFromString("abc"), "-", "--abc--"},
		{"odd padding", 8, TextFromString("abc"), "-", "--abc---"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignCenter(tt.width, tt.text, tt.filler)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestAlignLeft(t *testing.T) {
	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		filler   string
		expected string
	}{
		{"already fits", 3, TextFromString("abc"), ".", "abc"},
		{"padding", 6, TextFromString("abc"), ".", "...abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignLeft(tt.width, tt.text, tt.filler)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestAlignRight(t *testing.T) {
	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		filler   string
		expected string
	}{
		{"already fits", 3, TextFromString("abc"), ".", "abc"},
		{"padding", 6, TextFromString("abc"), ".", "abc..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AlignRight(tt.width, tt.text, tt.filler)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestPatternLeft(t *testing.T) {
	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		expected string
	}{
		{"already fits", 3, TextFromString("abc"), "abc"},
		{"single rune", 5, TextFromString("*"), "*****"},
		{"multiple runes exact", 6, TextFromString("ab"), "ababab"},
		{"multiple runes with remainder", 5, TextFromString("ab"), "babab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PatternLeft(tt.width, tt.text)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestPatternRight(t *testing.T) {
	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		expected string
	}{
		{"already fits", 3, TextFromString("abc"), "abc"},
		{"single rune", 5, TextFromString("*"), "*****"},
		{"multiple runes exact", 6, TextFromString("ab"), "ababab"},
		{"multiple runes with remainder", 5, TextFromString("ab"), "ababa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PatternRight(tt.width, tt.text)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestExtendLeft(t *testing.T) {
	result := ExtendLeft(5, TextFromString("abc"), ".")
	assert.Equal(t, result, ".....abc")
}

func TestExtendRight(t *testing.T) {
	result := ExtendRight(5, TextFromString("abc"), ".")
	assert.Equal(t, result, "abc.....")
}

func TestTruncateRight(t *testing.T) {
	ellipsis := Ellipsis{
		Data:  ".",
		Count: 3,
	}

	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		expected string
	}{
		{"empty text", 3, TextFromString(""), ""},
		{"fits", 10, TextFromString("abcdef"), "abcdef"},
		{"trim with ellipsis", 4, TextFromString("abcdef"), "a..."},
		{"ellipsis larger than width", 2, TextFromString("abcdef"), "ab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateRight(tt.width, tt.text, ellipsis)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestTruncateLeft(t *testing.T) {
	ellipsis := Ellipsis{
		Data:  ".",
		Count: 3,
	}

	tests := []struct {
		name     string
		width    winsize.Cols
		text     Text
		expected string
	}{
		{"empty text", 3, TextFromString(""), ""},
		{"fits", 10, TextFromString("abcdef"), "abcdef"},
		{"trim with ellipsis", 4, TextFromString("abcdef"), "...f"},
		{"ellipsis cannot fit", 2, TextFromString("abcdef"), "ef"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateLeft(tt.width, tt.text, ellipsis)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestNumberToAlpha(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"invalid zero", 0, "?"},
		{"invalid negative", -5, "?"},
		{"single letter a", 1, "a"},
		{"single letter b", 2, "b"},
		{"single letter z", 26, "z"},
		{"double letter aa", 27, "aa"},
		{"double letter ab", 28, "ab"},
		{"double letter az", 52, "az"},
		{"double letter ba", 53, "ba"},
		{"triple letter aaa", 703, "aaa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NumberToAlpha(tt.input)
			assert.Equal(t, result, tt.expected)
		})
	}
}
