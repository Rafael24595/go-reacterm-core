package splitter

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/delta"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
)

func SplitLine(line line.Line) ([]layout.Word, []layout.Frag) {
	return SplitLineWith(SplitFragByWords, line)
}

func SplitLineWith(splitter Frag, line line.Line) ([]layout.Word, []layout.Frag) {
	frgs := delta.New()

	for frg := range line.All() {
		if !frg.Atom().HasAny(atom.Wrap) && !frag.IsStructural(frg) {
			result := splitter(frg)
			frgs.Merge(result)
			continue
		}

		frgs.BoundAtEnd()
		frgs.AddFrag(*layout.NewFrag(&frg))
	}

	wrds := buildWords(frgs)
	return wrds, frgs.Frags
}

func SplitFragByWords(frg frag.Frag) delta.Delta {
	var sb strings.Builder

	var hasState bool
	var lastSpace bool

	var startsWithSpace bool

	frgs := delta.New()

	flush := func(src frag.Frag) {
		if sb.Len() == 0 {
			return
		}

		frg := frag.NewBuilder().
			AddText(sb.String()).
			WithMeta(src).
			Frag()

		wrd := layout.NewFrag(&frg)

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

func buildWords(frgs delta.Delta) []layout.Word {
	wrds := make([]layout.Word, 0, len(frgs.Bounds)+1)

	var wordStart uint32
	for _, p := range frgs.Bounds {
		wrds = append(
			wrds, *layout.New(wordStart, p),
		)

		wordStart = p
	}

	lenFrags := frgs.Size()
	if wordStart != lenFrags {
		wrds = append(
			wrds, *layout.New(wordStart, lenFrags),
		)
	}

	return wrds
}
