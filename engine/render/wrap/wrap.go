package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func NormalizeLines(lines ...line.Line) []LayoutLine {
	return normalizeLines(false, lines...)
}

func NormalizeLinesWithOrder(lines ...line.Line) []LayoutLine {
	return normalizeLines(true, lines...)
}

func normalizeLines(order bool, lines ...line.Line) []LayoutLine {
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
	for i, lne := range lines {
		if line.FragsMeasure(size.Cols, lne.Source) != 0 {
			continue
		}

		lastFrag := frag.Frag{}

		lneSize := lne.Source.Size()
		if lneSize > 0 {
			lastFrag = lne.Source.GetFrag(lneSize - 1)
		}

		frag := frag.NewBuilder().
			AddText(placeholder).
			WithMeta(&lastFrag).
			Frag()

		lines[i].Source = line.BuilderFromLine(lines[i].Source).
			PushFrags(frag).
			Line()

		lines[i].pushFrags(frag)
	}

	return lines
}

func Line(cols winsize.Cols, lne *line.Line) []line.Line {
	return wrapLine(cols, *lne, make([]line.Line, 0))
}

func Lines(cols winsize.Cols, lines ...line.Line) []line.Line {
	result := make([]line.Line, 0)

	for _, line := range lines {
		result = wrapLine(cols, line, result)
	}

	return result
}

func wrapLine(cols winsize.Cols, line line.Line, dst []line.Line) []line.Line {
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

func NextLine(cols winsize.Cols, lines []LayoutLine) (*line.Line, []LayoutLine) {
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

func wrapOnce(cols winsize.Cols, lne *LayoutLine) (*line.Line, *LayoutLine) {
	size := lne.Source.Size()

	cursor := line.NewBuilder(int(size)).
		WithMeta(lne.Source)

	remaining := cols
	currentWidth := winsize.Cols(0)

	wordIdx := 0

	for ; wordIdx < len(lne.words); wordIdx++ {
		wordMeasure := lne.measure(wordIdx, cols)

		if wordMeasure <= remaining {
			cursor.Text = appendFrags(
				cursor.Text, lne.findFrags(wordIdx),
			)

			remaining = remaining.Sub(wordMeasure)
			currentWidth += wordMeasure

			continue
		}

		if shouldWrap(lne, wordIdx, currentWidth) {
			break
		}

		if ok := lne.splitWord(
			wordIdx,
			cols,
			remaining,
		); ok {
			cursor.Text = appendFrags(
				cursor.Text, lne.findFrags(wordIdx),
			)
		}

		wordIdx++

		break
	}

	if wordIdx >= len(lne.words) {
		return cursor.LinePtr(), nil
	}

	rest := &LayoutLine{
		Source: lne.Source,
		frags:  lne.frags,
		words:  lne.words[wordIdx:],
	}

	return cursor.LinePtr(), rest
}

func shouldWrap(line *LayoutLine, wordIdx int, currentWidth winsize.Cols) bool {
	if line.hasAtom(wordIdx, atom.Break) {
		return false
	}

	return currentWidth > 0
}

func splitLineFeeds(lne *line.Line, order bool) []line.Line {
	result := make([]line.Line, 0)

	index := uint16(1)
	if lne.GetOrder() != 0 {
		index = lne.GetOrder()
	}

	builder := orderedBuilder(*lne, index, order)

	for frg := range lne.Frags() {
		if !strings.ContainsAny(frg.Text(), "\n\r") {
			builder.PushFrags(frg)
			continue
		}

		normalizedText := runes.NormalizeLineFeed(frg.Text())

		parts := strings.Split(normalizedText, "\n")
		for i, part := range parts {
			if part != "" {
				frgBuilder := frag.NewBuilder().
					AddText(part).
					WithMeta(&frg)

				builder.PushBuilder(frgBuilder)
			}

			if i >= len(parts)-1 {
				continue
			}

			result = append(result, builder.Line())
			index += 1

			builder = orderedBuilder(*lne, index, order)
		}
	}

	result = append(result, builder.Line())

	return result
}

func orderedBuilder(lne line.Line, index uint16, ordered bool) *line.Builder {
	current := line.NewBuilder().
		WithMeta(lne)

	if ordered {
		current.SetOrder(index)
	}

	return current
}

func splitFragAt(frg *wordFrag, cols winsize.Cols) (*wordFrag, *wordFrag) {
	if cols <= 0 {
		newFrag := frag.NewBuilder().
			WithMeta(frg.Base).
			Frag()

		return newWordFrag(&newFrag), frg
	}

	byteIndex, canBreak := runes.RuneIndexToByteIndex(frg.Base.Text(), offset.Offset(cols))
	if !canBreak || int(byteIndex) >= len(frg.Base.Text()) {
		return frg, nil
	}

	taken := frag.NewBuilder().
		AddText(frg.Base.Text()[:byteIndex]).
		WithMeta(frg.Base).
		Frag()

	rest := frag.NewBuilder().
		AddText(frg.Base.Text()[byteIndex:]).
		WithMeta(frg.Base).
		Frag()

	return newWordFrag(&taken), newWordFrag(&rest)
}
