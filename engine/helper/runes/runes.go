package runes

import (
	"slices"
	"strings"
	"unicode/utf8"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/ascii"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type measurable interface {
	~int | winsize.Cols | offset.Offset
}

var lineNormalizer = strings.NewReplacer("\r\n", "\n", "\r", "\n")

// WordDelimiters defines the set of runes that are considered as word boundaries for text navigation.
var WordDelimiters = []RuneRule{
	{
		Rune: ' ',
		Skip: false,
	},
	{
		Rune: '.',
		Skip: true,
	},
	{
		Rune: ',',
		Skip: true,
	},
	{
		Rune: ascii.ENTER_LF,
		Skip: true,
	},
}

// LineDelimiters defines the set of runes that are considered as line boundaries for text navigation.
var LineDelimiters = []RuneRule{
	{
		Rune: ascii.ENTER_LF,
		Skip: false,
	},
}

type RuneRule struct {
	// Rune is the target rune to search for.
	Rune rune
	// Skip indicates whether the rune should be skipped during searches.
	Skip bool
}

// NormalizeLineFeed converts \r\n and \r endings into \n.
func NormalizeLineFeed(text string) string {
	return lineNormalizer.Replace(text)
}

// Insert inserts a slice of runes in a buffer copy at a given offset position.
func Insert(
	buffer []rune,
	insert []rune,
	position offset.Offset,
) []rune {
	bufLen := offset.Offset(len(buffer))
	insLen := offset.Offset(len(insert))

	assert.False(
		bufLen < position,
		"position %d is out of bounds for buffer length %d", position, bufLen,
	)

	newBuffer := make([]rune, bufLen+insLen)

	copy(newBuffer[:position], buffer[:position])
	copy(newBuffer[position:], insert)
	copy(newBuffer[position+insLen:], buffer[position:])

	return newBuffer
}

// Replace replaces the slice range [start, end) with insert runes in a buffer copy.
func Replace(
	buffer []rune,
	insert []rune,
	start, end offset.Offset,
) []rune {
	if start == end {
		return Insert(buffer, insert, start)
	}

	bufLen := offset.Offset(len(buffer))
	insLen := offset.Offset(len(insert))

	assert.False(
		bufLen < end,
		"range[%d - %d] is greater than slice length %d", start, end, bufLen,
	)

	selLen := end.Sub(start)
	remLen := bufLen.Sub(selLen)

	newBuffer := make([]rune, remLen+insLen)

	copy(newBuffer[:start], buffer[:start])
	copy(newBuffer[start:], insert)
	copy(newBuffer[start+insLen:], buffer[end:])

	return newBuffer
}

// NormalizeBuffer ensures buffer has at least min capacity/length by padding with zero-value runes.
func NormalizeBuffer(buffer []rune, min uint) []rune {
	bufferLen := uint(len(buffer))
	if bufferLen >= min {
		return buffer
	}

	needed := math.SubClampZero(min, bufferLen)
	padding := make([]rune, needed)

	return append(buffer, padding...)
}

// BackwardIndexWithLimit navigates backward with an offset adjustment.
func BackwardIndexWithLimit[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	return BackwardIndex(buffer, definition, index) + 1
}

// BackwardIndex searches backward from index until matching a RuneDefinition.
func BackwardIndex[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	newIndex := fixedBackwardIndex(buffer, definition, index)

	j := math.SubClampZero(newIndex, 1)
	for {
		for i := range definition {
			if definition[i].Rune == buffer[j] {
				return j + 1
			}
		}

		if j == 0 {
			break
		}

		j = math.SubClampZero(j, 1)
	}

	return 0
}

func fixedBackwardIndex[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	newIndex := math.SubClampZero(index, 1)

	for newIndex > 0 {
		for i := range definition {
			v := definition[i]
			if v.Rune != buffer[newIndex] || !v.Skip {
				return newIndex
			}
		}

		newIndex--
	}

	return newIndex
}

// ForwardIndexWithLimit checks current index boundary before searching forward.
func ForwardIndexWithLimit[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	if index < T(len(buffer)) {
		for _, v := range definition {
			if buffer[index] == v.Rune {
				return index
			}
		}
	}

	return ForwardIndex(buffer, definition, index)
}

// ForwardIndex searches forward from index until matching a RuneDefinition.
func ForwardIndex[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	newIndex := fixForwardIndex(buffer, definition, index)

	for newIndex < T(len(buffer)) {
		for _, v := range definition {
			if v.Rune == buffer[newIndex] {
				return T(newIndex)
			}
		}

		newIndex++
	}

	return newIndex
}

func fixForwardIndex[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	bufferLen := T(len(buffer))
	newIndex := min(index+1, bufferLen)

	for newIndex < bufferLen {
		for _, v := range definition {
			if v.Rune != buffer[newIndex] || !v.Skip {
				return newIndex
			}
		}

		newIndex++
	}

	return newIndex
}

// JoinReverse concatenates a slice of strings in reverse order.
func JoinReverse(buffer []string) string {
	var sb strings.Builder
	for i := len(buffer) - 1; i >= 0; i-- {
		sb.WriteString(buffer[i])
	}
	return sb.String()
}

// RuneIndexToByteIndex translates a 0-based rune index into a byte index within string text.
func RuneIndexToByteIndex(text string, index offset.Offset) (offset.Offset, bool) {
	if index == 0 {
		return 0, true
	}

	count := offset.Offset(0)
	for i := range text {
		if count == index {
			return offset.Offset(i), true
		}
		count++
	}

	if count == index {
		return offset.Offset(len(text)), true
	}

	return 0, false
}

// SanitizeRunes removes null (0) runes from the slice in-place.
func SanitizeRunes(runes []rune) []rune {
	if !slices.Contains(runes, 0) {
		return runes
	}

	last := 0
	for _, r := range runes {
		if r != 0 {
			runes[last] = r
			last++
		}
	}

	return runes[:last]
}

// Measure returns the character count of a string as type T.
func Measure[T measurable](text string) T {
	return T(utf8.RuneCountInString(text))
}

// MeasureRunes returns the length of a rune slice as type T.
func MeasureRunes[T measurable](runes []rune) T {
	return T(len(runes))
}

// MeasureCols returns the character count of a string as winsize.Cols.
func MeasureCols(text string) winsize.Cols {
	return Measure[winsize.Cols](text)
}

// MeasureColsRunes returns the length of a rune slice as winsize.Cols.
func MeasureColsRunes(runes []rune) winsize.Cols {
	return MeasureRunes[winsize.Cols](runes)
}

// MeasureOffset returns the character count of a string as offset.Offset.
func MeasureOffset(text string) offset.Offset {
	return Measure[offset.Offset](text)
}

// MeasureOffsetRunes returns the length of a rune slice as offset.Offset.
func MeasureOffsetRunes(runes []rune) offset.Offset {
	return MeasureRunes[offset.Offset](runes)
}
