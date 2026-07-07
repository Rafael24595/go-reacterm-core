package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func NormalizeLines(lines ...text.Line) []LayoutLine {
	return normalizeLines(false, lines...)
}

func NormalizeLinesWithOrder(lines ...text.Line) []LayoutLine {
	return normalizeLines(true, lines...)
}

func normalizeLines(order bool, lines ...text.Line) []LayoutLine {
	buffer := make([]LayoutLine, 0, len(lines))

	for _, line := range lines {
		normalizedLF := splitLineFeeds(&line, order)

		for _, n := range normalizedLF {
			words := splitLineWords(&n)
			newLayoutLine := NewLayoutLine(&n, words...)
			buffer = append(buffer, *newLayoutLine)
		}
	}

	return buffer
}

func MaterializeEmpty(
	size winsize.Winsize,
	placeholder string,
	lines ...LayoutLine,
) []LayoutLine {
	for i, line := range lines {
		if text.FragmentMeasure(size.Cols, line.Source.Text...) != 0 {
			continue
		}

		lastFrag := text.Fragment{}
		if len(line.Source.Text) > 0 {
			lastFrag = line.Source.Text[len(line.Source.Text)-1]
		}

		fragment := *text.NewFragment(placeholder).
			CopyMeta(&lastFrag)

		lines[i].Source.PushFragments(fragment)
		lines[i].Words = append(line.Words, *newWord(fragment))
	}

	return lines
}

func Line(cols winsize.Cols, line *text.Line) []text.Line {
	return wrapLine(cols, *line, make([]text.Line, 0))
}

func Lines(cols winsize.Cols, lines ...text.Line) []text.Line {
	result := make([]text.Line, 0)

	for _, line := range lines {
		result = wrapLine(cols, line, result)
	}

	return result
}

func wrapLine(cols winsize.Cols, line text.Line, dst []text.Line) []text.Line {
	words := splitLineWords(&line)
	layout := NewLayoutLine(&line, words...)

	current := layout

	for current != nil {
		head, rest := wrapOnce(cols, *current)
        dst = append(dst, *head)
        current = rest
	}

	return dst
}

func NextLine(cols winsize.Cols, lines []LayoutLine) (*text.Line, []LayoutLine) {
	if cols == 0 || len(lines) == 0 {
		return nil, make([]LayoutLine, 0)
	}

	current := lines[0]
	remain := lines[1:]

	result, rest := wrapOnce(cols, current)
	if rest != nil {
		remain = append([]LayoutLine{*rest}, remain...)
	}

	return result, remain
}

func wrapOnce(cols winsize.Cols, line LayoutLine) (*text.Line, *LayoutLine) {
	cursor := text.LineFromMeta(line.Source, len(line.Source.Text))

	remaining := cols
	currentWidth := winsize.Cols(0)

	words := line.Words

	for len(words) > 0 {
		focus := &words[0]

		wordMeasure := focus.Measure(cols)

		if wordMeasure <= remaining {
			cursor.Text = appendFragments(
				cursor.Text, focus.Text...,
			)

			remaining = remaining.Sub(wordMeasure)
			currentWidth += wordMeasure
			words = words[1:]

			continue
		}

		if shouldWrap(*focus, currentWidth) {
			break
		}

		words = words[1:]

		newWord, restWord := splitLongWord(*focus, cols, remaining)
		if newWord != nil {
			cursor.Text = appendFragments(
				cursor.Text, newWord.Text...,
			)
		}

		if restWord != nil {
			words = append([]word{*restWord}, words...)
		}

		break
	}

	if len(words) == 0 {
		return cursor, nil
	}

	rest := NewLayoutLine(line.Source, words...)

	return cursor, rest
}

func shouldWrap(word word, currentWidth winsize.Cols) bool {
	if word.HasAtom(atom.Break) {
		return false
	}

	return currentWidth > 0
}

func splitLineFeeds(line *text.Line, order bool) []text.Line {
	result := make([]text.Line, 0)

	index := uint16(1)
	if line.Order != 0 {
		index = line.Order
	}

	current := text.LineFromMeta(line)
	if order {
		current.SetOrder(index)
	}

	for _, frag := range line.Text {
		if !strings.ContainsAny(frag.Text, "\n\r") {
			current.PushFragments(frag)
			continue
		}

		normalizedText := runes.NormalizeLineFeed(frag.Text)

		parts := strings.Split(normalizedText, "\n")
		for i, part := range parts {
			if part != "" {
				current.PushFragments(
					*text.NewFragment(part).CopyMeta(&frag),
				)
			}

			if i >= len(parts)-1 {
				continue
			}

			result = append(result, *current)
			index += 1

			current = text.LineFromMeta(line)
			if order {
				current.SetOrder(index)
			}
		}
	}

	result = append(result, *current)

	return result
}

func splitFragmentAt(frag *wordFrag, cols winsize.Cols) (*wordFrag, *wordFrag) {
	if cols <= 0 {
		newFrag := text.EmptyFragment().
			CopyMeta(frag.Base)
		return newWordFrag(newFrag), frag
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(frag.Base.Text, offset.Offset(cols))
	if !canBreak || int(byteIndex) >= len(frag.Base.Text)  {
		return frag, nil
	}

	taken := text.NewFragment(frag.Base.Text[:byteIndex]).
		CopyMeta(frag.Base)

	rest := text.NewFragment(frag.Base.Text[byteIndex:]).
		CopyMeta(frag.Base)

	return newWordFrag(taken), newWordFrag(rest)
}
