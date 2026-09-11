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

var LineDelimiters = []RuneRule{
	{
		Rune: ascii.ENTER_LF,
		Skip: false,
	},
}

type RuneRule struct {
	Rune rune
	Skip bool
}

func NormalizeLineFeed(text string) string {
	return lineNormalizer.Replace(text)
}

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

func NormalizeBuffer(buffer []rune, min uint) []rune {
	bufferLen := uint(len(buffer))
	if bufferLen >= min {
		return buffer
	}

	needed := math.SubClampZero(min, bufferLen)
	padding := make([]rune, needed)

	return append(buffer, padding...)
}

func BackwardIndexWithLimit[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	return BackwardIndex(buffer, definition, index) + 1
}

func BackwardIndex[T math.Number](
	buffer []rune,
	definition []RuneRule,
	index T,
) T {
	newIndex := fixedBackwardIndex(buffer, definition, index)

	j := math.SubClampZero(newIndex, 1)
	for {
		for _, v := range definition {
			if v.Rune == buffer[j] {
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
		for _, v := range definition {
			if v.Rune != buffer[newIndex] || !v.Skip {
				return newIndex
			}
		}

		newIndex--
	}

	return newIndex
}

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

func JoinReverse(buffer []string) string {
	var sb strings.Builder
	for i := len(buffer) - 1; i >= 0; i-- {
		sb.WriteString(buffer[i])
	}
	return sb.String()
}

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

func Measure[T measurable](text string) T {
	return T(utf8.RuneCountInString(text))
}

func MeasureRunes[T measurable](runes []rune) T {
	return T(len(runes))
}

func MeasureCols(text string) winsize.Cols {
	return Measure[winsize.Cols](text)
}

func MeasureColsRunes(runes []rune) winsize.Cols {
	return MeasureRunes[winsize.Cols](runes)
}

func MeasureOffset(text string) offset.Offset {
	return Measure[offset.Offset](text)
}

func MeasureOffsetRunes(runes []rune) offset.Offset {
	return MeasureRunes[offset.Offset](runes)
}
