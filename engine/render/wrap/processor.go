package wrap

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Processor func(order bool, lne line.Line) []line.Line

func LineFeedProcessor(order bool, lne line.Line) []line.Line {
	result := make([]line.Line, 0)

	index := uint16(1)
	if lne.Order() != 0 {
		index = lne.Order()
	}

	builder := orderedBuilder(lne, index, order)

	for frg := range lne.All() {
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
					WithMeta(frg)

				builder.PushBuilder(frgBuilder)
			}

			if i >= len(parts)-1 {
				continue
			}

			result = append(result, builder.Line())
			index += 1

			builder = orderedBuilder(lne, index, order)
		}
	}

	result = append(result, builder.Line())

	return result
}
