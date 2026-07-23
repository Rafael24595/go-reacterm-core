package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type word struct {
	start    uint32
	end      uint32
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func splitLineWords(line line.Line) ([]word, []wordFrag) {
	words := make([]word, 0, line.Size())
	frags := make([]wordFrag, 0, line.Size())

	var sb strings.Builder
	var lastSpace bool
	var hasState bool

	wordStart := 0

	flushFrag := func(frg frag.Frag) {
		if sb.Len() == 0 {
			return
		}

		f := frag.NewBuilder().
			AddText(sb.String()).
			WithMeta(&frg).
			Frag()

		frags = append(frags, *newWordFrag(&f))
		sb.Reset()
	}

	flushWord := func() {
		if wordStart == len(frags) {
			return
		}

		words = append(words, word{
			start: uint32(wordStart),
			end:   uint32(len(frags)),
		})

		wordStart = len(frags)
	}

	for frg := range line.All() {
		if frg.Atom().HasAny(atom.Wrap) || frag.IsStructural(frg) {
			flushFrag(frg)
			flushWord()

			frags = append(frags, *newWordFrag(&frg))

			words = append(words, word{
				start: uint32(wordStart),
				end:   uint32(wordStart + 1),
			})

			wordStart = len(frags)
			hasState = false

			continue
		}

		for _, r := range frg.Text() {
			isSpace := unicode.IsSpace(r)

			if hasState && isSpace != lastSpace {
				flushFrag(frg)
				flushWord()
			}

			lastSpace = isSpace
			hasState = true

			sb.WriteRune(r)
		}

		flushFrag(frg)
	}

	flushWord()

	return words, frags
}
