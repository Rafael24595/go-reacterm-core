package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type LineSplitter func(lne line.Line) ([]word, []wordFrag)
type FragSplitter func(frg frag.Frag) Delta

func SplitLine(line line.Line) ([]word, []wordFrag) {
	return SplitLineWith(SplitFragByWords, line)
}

func SplitLineWith(splitter FragSplitter, line line.Line) ([]word, []wordFrag) {
	frgs := NewDelta()

	for frg := range line.All() {
		if !frg.Atom().HasAny(atom.Wrap) && !frag.IsStructural(frg) {
			result := splitter(frg)
			frgs.Merge(result)
			continue
		}

		frgs.BoundAtEnd()
		frgs.AddFrag(*newWordFrag(&frg))
	}

	wrds := buildWords(frgs)
	return wrds, frgs.Frags
}

func SplitFragByWords(frg frag.Frag) Delta {
	var sb strings.Builder

	var hasState bool
	var lastSpace bool

	var startsWithSpace bool

	frgs := NewDelta()

	flush := func(src frag.Frag) {
		if sb.Len() == 0 {
			return
		}

		frg := frag.NewBuilder().
			AddText(sb.String()).
			WithMeta(src).
			Frag()

		wrd := newWordFrag(&frg)

		frgs.AddFrag(*wrd)

		sb.Reset()
	}

	for _, r := range frg.Text() {
		isSpace := unicode.IsSpace(r)

		if !hasState {
			startsWithSpace = isSpace
		} else if isSpace != lastSpace {
			flush(frg)
			frgs.BoundAtEnd()
		}

		lastSpace = isSpace
		hasState = true

		sb.WriteRune(r)
	}

	flush(frg)

	frgs.LeftEdge = startsWithSpace
	frgs.RightEdge = lastSpace

	return frgs
}

func buildWords(frgs Delta) []word {
	wrds := make([]word, 0, len(frgs.Bounds)+1)

	var wordStart uint32
	for _, p := range frgs.Bounds {
		wrds = append(
			wrds, *newWord(wordStart, p),
		)

		wordStart = p
	}

	lenFrags := frgs.Size()
	if wordStart != lenFrags {
		wrds = append(
			wrds, *newWord(wordStart, lenFrags),
		)
	}

	return wrds
}
