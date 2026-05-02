package helper

import (
	"fmt"
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

type TextLayoutOpts struct {
	Text        string
	LogicalSize winsize.Cols
}

func fixDirectionOps(text string, opts TextLayoutOpts) TextLayoutOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	if opts.Text == "" {
		opts.Text = marker.DefaultPaddingText
	}

	return opts
}

type LogicalSizeOpts struct {
	LogicalSize winsize.Cols
}

func fixLogicalSizeOpts(text string, opts LogicalSizeOpts) LogicalSizeOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	return opts
}

type TextTrimOpts struct {
	LogicalSize  winsize.Cols
	EllipsisText string
	EllipsisSize winsize.Cols
}

func fixTextTrimOpts(text string, opts TextTrimOpts) TextTrimOpts {
	if opts.LogicalSize == 0 {
		opts.LogicalSize = runes.Measure(text)
	}

	if opts.EllipsisSize == 0 {
		opts.EllipsisSize = marker.DefaultElipsisSize
	}

	return opts
}

func Center(item any, width winsize.Cols) string {
	return CenterWithOpts(item, width, TextLayoutOpts{})
}

func CenterWithOpts(item any, width winsize.Cols, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width.Clamp(opts.LogicalSize)

	left := int(padding / 2)
	right := int(padding) - left

	return strings.Repeat(opts.Text, left) + text + strings.Repeat(opts.Text, right)
}

func Left(item any, width winsize.Cols) string {
	return LeftWithOpts(item, width, TextLayoutOpts{})
}

func LeftWithOpts(item any, width winsize.Cols, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width.Clamp(opts.LogicalSize)

	return strings.Repeat(opts.Text, int(padding)) + text
}

func Right(item any, width winsize.Cols) string {
	return RightWithOpts(item, width, TextLayoutOpts{})
}

func RightWithOpts(item any, width winsize.Cols, opts TextLayoutOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixDirectionOps(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	padding := width.Clamp(opts.LogicalSize)

	return text + strings.Repeat(opts.Text, int(padding))
}

func FillLeft(item any, width winsize.Cols) string {
	return FillLeftWithOpts(item, width, LogicalSizeOpts{})
}

func FillLeftWithOpts(item any, width winsize.Cols, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixLogicalSizeOpts(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	if text == "" {
		text = marker.DefaultPaddingText
	}

	textLen := runes.Measure(text)

	fix := ""
	if rest := width % textLen; rest != 0 {
		fix = text[rest:]
	}

	width = width / textLen

	return fix + strings.Repeat(text, int(width))
}

func FillRight(item any, width winsize.Cols) string {
	return FillRightWithOpts(item, width, LogicalSizeOpts{})
}

func FillRightWithOpts(item any, width winsize.Cols, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)

	opts = fixLogicalSizeOpts(text, opts)
	if opts.LogicalSize >= width {
		return text
	}

	if text == "" {
		text = marker.DefaultPaddingText
	}

	textLen := runes.Measure(text)

	fix := ""
	if rest := width % textLen; rest != 0 {
		fix = text[:rest]
	}

	width = width / textLen

	return strings.Repeat(text, int(width)) + fix
}

func RepeatLeft(item any, runes string, width winsize.Cols) string {
	return RepeatLeftWithOpts(item, runes, width, LogicalSizeOpts{})
}

func RepeatLeftWithOpts(item any, runes string, width winsize.Cols, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)
	return FillLeftWithOpts(runes, width, opts) + text
}

func RepeatRight(item any, runes string, width winsize.Cols) string {
	return RepeatRightWithOpts(item, runes, width, LogicalSizeOpts{})
}

func RepeatRightWithOpts(item any, runes string, width winsize.Cols, opts LogicalSizeOpts) string {
	text := fmt.Sprintf("%v", item)
	return text + FillRightWithOpts(runes, width, opts)
}

func TrimLeft(data string, width winsize.Cols, opts TextTrimOpts) string {
	if data == "" {
		return data
	}

	opts = fixTextTrimOpts(data, opts)

	elipSize := runes.Measure(opts.EllipsisText) * opts.EllipsisSize

	realSize := runes.Measure(data)
	if width >= opts.LogicalSize || width > realSize {
		return data
	}

	width = realSize.Clamp(width)

	if elipSize+width >= realSize {
		index, _ := runes.RuneIndexToByteIndex(data, offset.Offset(width))
		return data[index:]
	}

	elipTotal := strings.Repeat(opts.EllipsisText, int(opts.EllipsisSize))

	index, _ := runes.RuneIndexToByteIndex(data, offset.Offset(width+elipSize))
	return elipTotal + data[index:]
}

func TrimRight(data string, width winsize.Cols, opts TextTrimOpts) string {
	if data == "" {
		return data
	}

	opts = fixTextTrimOpts(data, opts)

	elipSize := runes.Measure(opts.EllipsisText) * opts.EllipsisSize

	realSize := runes.Measure(data)
	if width >= opts.LogicalSize || width > realSize {
		return data
	}

	if elipSize > width {
		index, _ := runes.RuneIndexToByteIndex(data, offset.Offset(width))
		return data[:index]
	}

	elipTotal := strings.Repeat(opts.EllipsisText, int(opts.EllipsisSize))

	size := width.Clamp(elipSize)
	index, _ := runes.RuneIndexToByteIndex(data, offset.Offset(size))

	return data[:index] + elipTotal
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
