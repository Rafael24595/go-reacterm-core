package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
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
			words, frags := splitLineWords(&n)
			layout := NewLayoutLine(&n, words, frags)
			buffer = append(buffer, *layout)
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
		if frag.Measure(size.Cols, line.Source.Text...) != 0 {
			continue
		}

		lastFrag := frag.Frag{}
		if len(line.Source.Text) > 0 {
			lastFrag = line.Source.Text[len(line.Source.Text)-1]
		}

		frag := *frag.New(placeholder).
			CopyMeta(&lastFrag)

		lines[i].Source.PushFrags(frag)
		lines[i].pushFrags(frag)
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
	words, frags := splitLineWords(&line)
	layout := NewLayoutLine(&line, words, frags)

	current := layout

	for current != nil {
		head, rest := wrapOnce(cols, current)
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

	result, rest := wrapOnce(cols, &current)
	if rest != nil {
		remain = append([]LayoutLine{*rest}, remain...)
	}

	return result, remain
}

func wrapOnce(cols winsize.Cols, line *LayoutLine) (*text.Line, *LayoutLine) {
	cursor := text.LineFromMeta(line.Source, len(line.Source.Text))

	remaining := cols
	currentWidth := winsize.Cols(0)

	wordIdx := 0

	for ; wordIdx < len(line.words); wordIdx++ {
		wordMeasure := line.measure(wordIdx, cols)

		if wordMeasure <= remaining {
			cursor.Text = appendFrags(
				cursor.Text,
				line.findFrags(wordIdx),
			)

			remaining = remaining.Sub(wordMeasure)
			currentWidth += wordMeasure

			continue
		}

		if shouldWrap(line, wordIdx, currentWidth) {
			break
		}

		if ok := line.splitWord(
			wordIdx,
			cols,
			remaining,
		); ok {
			cursor.Text = appendFrags(
				cursor.Text, line.findFrags(wordIdx),
			)
		}

		wordIdx++

		break
	}

	if wordIdx >= len(line.words) {
		return cursor, nil
	}

	rest := &LayoutLine{
		Source: line.Source,
		frags:  line.frags,
		words:  line.words[wordIdx:],
	}

	return cursor, rest
}

func shouldWrap(line *LayoutLine, wordIdx int, currentWidth winsize.Cols) bool {
	if line.hasAtom(wordIdx, atom.Break) {
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

	for _, frg := range line.Text {
		if !strings.ContainsAny(frg.Text, "\n\r") {
			current.PushFrags(frg)
			continue
		}

		normalizedText := runes.NormalizeLineFeed(frg.Text)

		parts := strings.Split(normalizedText, "\n")
		for i, part := range parts {
			if part != "" {
				current.PushFrags(
					*frag.New(part).CopyMeta(&frg),
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

func splitFragAt(frg *wordFrag, cols winsize.Cols) (*wordFrag, *wordFrag) {
	if cols <= 0 {
		newFrag := frag.Empty().
			CopyMeta(frg.Base)

		return newWordFrag(newFrag), frg
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(frg.Base.Text, offset.Offset(cols))
	if !canBreak || int(byteIndex) >= len(frg.Base.Text) {
		return frg, nil
	}

	taken := frag.New(frg.Base.Text[:byteIndex]).
		CopyMeta(frg.Base)

	rest := frag.New(frg.Base.Text[byteIndex:]).
		CopyMeta(frg.Base)

	return newWordFrag(taken), newWordFrag(rest)
}
