package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type word struct {
	Text []text.Fragment
}

func newWord(text ...text.Fragment) *word {
	return &word{
		Text: text,
	}
}

func (w *word) addText(text ...text.Fragment) *word {
	w.Text = append(w.Text, text...)
	return w
}

func splitLineWords(line *text.Line) []word {
	tokens := make([]word, 0, len(line.Text))
	frags := make([]text.Fragment, 0, 4)

	var sb strings.Builder
	var lastSpace bool
	var hasState bool

	flushFrag := func(frag text.Fragment) {
		if sb.Len() == 0 {
			return
		}

		f := text.NewFragment(sb.String()).
			CopyMeta(&frag)

		frags = append(frags, *f)

		sb.Reset()
	}

	flushWord := func() {
		if len(frags) == 0 {
			return
		}

		tokenFrags := make([]text.Fragment, len(frags))
		copy(tokenFrags, frags)

		tokens = append(tokens, word{
			Text: tokenFrags,
		})

		frags = frags[:0]
	}

	for _, frag := range line.Text {
		if frag.Atom.HasAny(style.AtmWrap) || text.IsStructuralFragment(frag) {
			flushFrag(frag)
			flushWord()

			tokens = append(tokens, word{
				Text: []text.Fragment{frag},
			})

			hasState = false
			continue
		}

		for _, r := range frag.Text {
			isSpace := unicode.IsSpace(r)

			if hasState && isSpace != lastSpace {
				flushFrag(frag)
				flushWord()
			}

			lastSpace = isSpace
			hasState = true

			sb.WriteRune(r)
		}

		flushFrag(frag)
	}

	flushWord()

	return tokens
}

func splitLongWord(
	word word,
	cols winsize.Cols,
	remaining winsize.Cols,
) (*word, *word) {
	if cols == 0 || remaining == 0 {
		return nil, &word
	}

	current := newWord()
	frags := word.Text

	for len(frags) > 0 {
		frag := frags[0]
		size := text.FragmentMeasure(cols, frag)

		if size <= remaining {
			current.Text = append(current.Text, frag)
			remaining = remaining.Clamp(size)
			frags = frags[1:]

			continue
		}

		takenFrag, restFrag := splitFragmentAt(&frag, remaining)

		current.Text = append(current.Text, *takenFrag)

		rest := newWord(*restFrag).
			addText(frags[1:]...)

		return current, rest
	}

	return current, nil
}

func wordsToLine(base text.Line, words ...word) *text.Line {
	line := text.EmptyLine().
		CopyMeta(&base)

	for _, w := range words {
		line.PushFragments(w.Text...)
	}

	return line
}
