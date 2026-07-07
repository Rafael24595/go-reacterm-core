package wrap

import (
	"strings"
	"unicode"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type word struct {
	start    uint32
	end      uint32
	measured bool
	cols     winsize.Cols
	measure  winsize.Cols
}

func splitLineWords(line *text.Line) ([]word, []wordFrag) {
	words := make([]word, 0, len(line.Text))
	frags := make([]wordFrag, 0, len(line.Text))

	var sb strings.Builder
	var lastSpace bool
	var hasState bool

	wordStart := 0

	flushFrag := func(frag text.Frag) {
		if sb.Len() == 0 {
			return
		}

		f := text.NewFrag(sb.String()).
			CopyMeta(&frag)

		frags = append(frags, *newWordFrag(f))
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

	for _, frag := range line.Text {
		if frag.Atom.HasAny(atom.Wrap) || text.IsStructuralFrag(frag) {
			flushFrag(frag)
			flushWord()

			frags = append(frags, *newWordFrag(&frag))

			words = append(words, word{
				start: uint32(wordStart),
				end:   uint32(wordStart + 1),
			})

			wordStart = len(frags)
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

	return words, frags
}
