package rule

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

var wrappers = map[rune]rune{
	'{': '}',
	'(': ')',
	'[': ']',
	'<': '>',
}

var trailingSpaceRunes = set.From(
	',',
	'.',
	';',
)

type Rule func(
	text []rune,
	start, end offset.Offset,
	buff []rune,
) ([]rune, bool)

var Full = []Rule{
	AppendSpaceAfter, WrapSelection,
}

func WrapSelection(
	input []rune,
	start, end offset.Offset,
	buffer []rune,
) ([]rune, bool) {
	if len(input) != 1 {
		return input, false
	}

	openRune := input[0]
	
	closeRune, ok := wrappers[openRune]
	if !ok {
		return input, false
	}

	buffLen := offset.Offset(len(buffer))
	if start > end || start > buffLen || end > buffLen {
		assert.Unreachable("Invalid offset values: start=%d, end=%d, buffer length=%d", start, end, buffLen)
		return input, false
	}

	selected := buffer[start:end]

	result := make([]rune, 0, len(selected)+2)
	result = append(result, openRune)
	result = append(result, selected...)
	result = append(result, closeRune)

	return result, true
}

func AppendSpaceAfter(
	input []rune,
	start, end offset.Offset,
	_ []rune,
) ([]rune, bool) {
	if len(input) != 1 {
		return input, false
	}

	inputRune := input[0]
	if !trailingSpaceRunes.Has(inputRune) {
		return input, false
	}

	result := []rune{input[0], ' '}
	return result, true
}
