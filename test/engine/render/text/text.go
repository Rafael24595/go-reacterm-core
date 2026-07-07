package text_test

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func FragsToString(frags []text.Frag) string {
	var s strings.Builder
	for _, f := range frags {
		s.WriteString(f.Text)
	}
	return s.String()
}

func LineToString(line *text.Line) string {
	buffer := make([]string, 0)
	for _, v := range line.Text {
		buffer = append(buffer, v.Text)
	}
	return strings.Join(buffer, "")
}
