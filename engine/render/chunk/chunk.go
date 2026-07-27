package chunk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

const DefaultChunk = 64

func Line(src line.Line, limit offset.Offset) line.Line {
	if limit == 0 {
		return src
	}

	if src.Measure() <= winsize.Cols(limit) {
		return src
	}

	builder := line.NewBuilder().
		WithMeta(src)

	for frag := range src.All() {
		split(builder, frag, limit)
	}

	return builder.Line()
}

func split(builder *line.Builder, frg frag.Frag, limit offset.Offset) {
	for {
		head, tail := splitAt(frg, limit)
		builder.PushFrags(head)

		if tail == nil {
			break
		}

		frg = *tail
	}
}

//TODO: Handle special atoms?
func splitAt(frg frag.Frag, limit offset.Offset) (frag.Frag, *frag.Frag) {
	text := frg.Text()

	byteIndex, canBreak := runes.RuneIndexToByteIndex(text, limit)
	if !canBreak || int(byteIndex) >= len(text) {
		return frg, nil
	}

	head := clone(frg, text[:byteIndex])
	tail := clone(frg, text[byteIndex:])

	return head, &tail
}

func clone(frg frag.Frag, text string) frag.Frag {
	return frag.NewBuilder().
		AddText(text).
		WithMeta(frg).
		Frag()
}
