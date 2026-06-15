package format

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

func AlignCenter(width winsize.Cols, text Text, filler string) string {
	if text.Size >= width {
		return text.Data
	}

	padding := width.Sub(text.Size)

	paddingLeft := int(padding / 2)
	paddingRight := int(padding) - paddingLeft

	left := strings.Repeat(filler, paddingLeft)
	right := strings.Repeat(filler, paddingRight)

	return left + text.Data + right
}

func AlignLeft(width winsize.Cols, text Text, filler string) string {
	if text.Size >= width {
		return text.Data
	}

	padding := width.Sub(text.Size)

	return strings.Repeat(filler, int(padding)) + text.Data
}

func AlignRight(width winsize.Cols, text Text, filler string) string {
	if text.Size >= width {
		return text.Data
	}

	padding := width.Sub(text.Size)

	return text.Data + strings.Repeat(filler, int(padding))
}

func PatternLeft(width winsize.Cols, text Text) string {
	if text.Size >= width {
		return text.Data
	}

	data := text.Data
	if text.Data == "" {
		data = marker.DefaultPaddingText
	}

	measure := runes.Measure(data)

	fix := ""
	if rest := width % measure; rest != 0 {
		fix = data[rest:]
	}

	width = width / measure

	return fix + strings.Repeat(data, int(width))
}

func PatternRight(width winsize.Cols, text Text) string {
	if text.Size >= width {
		return text.Data
	}

	data := text.Data
	if text.Data == "" {
		data = marker.DefaultPaddingText
	}

	measure := runes.Measure(data)

	fix := ""
	if rest := width % measure; rest != 0 {
		fix = data[:rest]
	}

	width = width / measure

	return strings.Repeat(data, int(width)) + fix
}

func ExtendLeft(width winsize.Cols, text Text, filler string) string {
	return PatternLeft(width, TextFromString(filler)) + text.Data
}

func ExtendRight(width winsize.Cols, text Text, filler string) string {
	return text.Data + PatternRight(width, TextFromString(filler))
}

func TruncateLeft(width winsize.Cols, text Text, ellipsis Ellipsis) string {
	if text.Data == "" {
		return text.Data
	}

	realSize := runes.Measure(text.Data)
	if width >= text.Size || width > realSize {
		return text.Data
	}

	width = realSize.Sub(width)
	ellipsisMeasure := ellipsis.measure()

	if ellipsisMeasure+width >= realSize {
		offset := offset.Offset(width)
		index, _ := runes.RuneIndexToByteIndex(text.Data, offset)
		return text.Data[index:]
	}

	ellipsisText := strings.Repeat(ellipsis.Data, int(ellipsis.Count))

	offset := offset.Offset(width + ellipsisMeasure)
	index, _ := runes.RuneIndexToByteIndex(text.Data, offset)

	return ellipsisText + text.Data[index:]
}

func TruncateRight(width winsize.Cols, text Text, ellipsis Ellipsis) string {
	if text.Data == "" {
		return text.Data
	}

	realSize := runes.Measure(text.Data)
	if width >= text.Size || width > realSize {
		return text.Data
	}

	ellipsisMeasure := ellipsis.measure()
	if ellipsisMeasure > width {
		offset := offset.Offset(width)
		index, _ := runes.RuneIndexToByteIndex(text.Data, offset)
		return text.Data[:index]
	}

	ellipsisText := strings.Repeat(ellipsis.Data, int(ellipsis.Count))

	size := width.Sub(ellipsisMeasure)
	index, _ := runes.RuneIndexToByteIndex(text.Data, offset.Offset(size))

	return text.Data[:index] + ellipsisText
}

func NumberToAlpha(n int) string {
	if n <= 0 {
		return "?"
	}

	result := ""

	for n > 0 {
		n--
		remainder := n % 26
		result = string(rune('a'+remainder)) + result
		n = n / 26
	}

	return result
}
