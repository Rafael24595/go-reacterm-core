package rule

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

var wrapperMap = map[rune]rune{
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
	text []rune,
	start, end offset.Offset,
	buff []rune,
) ([]rune, bool) {
	size := len(text)
	if size < 1 || size > 1 {
		return text, false
	}

	focus := text[0]

	close, ok := wrapperMap[focus]
	if !ok {
		return text, false
	}

	text = make([]rune, 0)
	text = append(text, focus)
	text = append(text, buff[start:end]...)
	text = append(text, close)

	return text, true
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
