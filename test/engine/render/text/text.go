package text_test

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func FragsToString(frags []frag.Frag) string {
	var s strings.Builder
	for _, f := range frags {
		s.WriteString(f.Text())
	}
	return s.String()
}

func LineToString(line line.Line) string {
	return FragsToString(line.GetText())
}
