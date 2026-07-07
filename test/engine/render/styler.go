package render_test

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Frags(styler *styler.Spec, size winsize.Winsize, frags []text.Frag) string {
	var buffer strings.Builder

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range frags {
		txt := format.NewText(
			f.Text,
			f.Size(),
		)
		spec := styler.Apply(f.Spec, lineSize, txt)

		fragSize := text.FragsMeasure(size.Cols, f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		buffer.WriteString(spec)
	}

	return buffer.String()
}
