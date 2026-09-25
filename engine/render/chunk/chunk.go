package chunk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// DefaultMeasure represents the default maximum fragment width in character cells.
const DefaultMeasure = 64

// Line partitions a line's fragments so that no individual fragment exceeds the specified offset limit.
func Line(src line.Line, limit offset.Offset) line.Line {
	if limit == 0 {
		return src
	}

	if src.Measure() <= winsize.Cols(limit) {
		return src
	}

	builder := line.NewBuilder().
		WithMeta(src)

	for frg := range src.All() {
		split(builder, *frg, limit)
	}

	return builder.Line()
}

func split(
	builder *line.Builder,
	frg frag.Frag,
	limit offset.Offset,
) {
	for {
		head, tail, hasTail := splitAt(frg, limit)
		builder.PushFrags(head)

		if !hasTail {
			break
		}

		frg = tail
	}
}

// TODO: Handle special atoms?
func splitAt(frg frag.Frag, limit offset.Offset) (frag.Frag, frag.Frag, bool) {
	text := frg.Text()

	byteIndex, canBreak := runes.RuneIndexToByteIndex(text, limit)
	if !canBreak || int(byteIndex) >= len(text) {
		return frg, frag.Frag{}, false
	}

	head := clone(frg, text[:byteIndex])
	tail := clone(frg, text[byteIndex:])

	return head, tail, true
}

func clone(frg frag.Frag, text string) frag.Frag {
	return frag.NewBuilder().
		AddText(text).
		WithMeta(frg).
		Frag()
}
