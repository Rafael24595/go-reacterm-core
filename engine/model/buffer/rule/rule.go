package rule

import "github.com/Rafael24595/go-reacterm-core/engine/model/offset"

var wrapperMap = map[rune]rune{
	'{': '}',
	'(': ')',
	'[': ']',
	'<': '>',
}

var runesRequiringTrailingSpace = []rune{
	',',
	'.',
	';',
}

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
	text []rune,
	start, end offset.Offset,
	_ []rune,
) ([]rune, bool) {
	size := len(text)
	if size < 1 || size > 1 {
		return text, false
	}

	focus := text[0]
	for _, r := range runesRequiringTrailingSpace {
		if focus != r {
			continue
		}

		text = append(text, ' ')
		return text, true
	}

	return text, false
}
