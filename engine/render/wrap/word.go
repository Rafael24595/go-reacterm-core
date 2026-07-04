package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type measureResolver func(winsize.Cols, ...text.Fragment) winsize.Cols

type word struct {
	Text     []text.Fragment
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func newWord(frags ...text.Fragment) *word {
	return &word{
		Text:     frags,
		measured: false,
		cols:     0,
		measure:  0,
	}
}

func (w *word) Measure(cols winsize.Cols) winsize.Cols {
	return w.measureWith(cols, text.FragmentMeasure)
}

func (w *word) measureWith(
	cols winsize.Cols,
	resolver measureResolver,
) winsize.Cols {
	if !w.measured || w.cols != cols {
		w.measure = resolver(cols, w.Text...)
		w.cols = cols
		w.measured = true
	}

	return w.measure
}

func splitLineWords(line *text.Line) []word {
	words := make([]word, 0, len(line.Text))
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

		words = append(words, word{
			Text: tokenFrags,
		})

		frags = frags[:0]
	}

	for _, frag := range line.Text {
		if frag.Atom.HasAny(atom.Wrap) || text.IsStructuralFragment(frag) {
			flushFrag(frag)
			flushWord()

			words = append(words, word{
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

	return words
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
			remaining = remaining.Sub(size)
			frags = frags[1:]

			continue
		}

		takenFrag, restFrag := splitFragmentAt(&frag, remaining)
		current.Text = append(current.Text, *takenFrag)

		rest := make([]text.Fragment, 0)
		if restFrag != nil {
			rest = append(rest, *restFrag)
		}

		rest = append(rest, frags[1:]...)
		if len(rest) == 0 {
			return current, nil
		}

		return current, newWord(rest...)
	}

	return current, nil
}
