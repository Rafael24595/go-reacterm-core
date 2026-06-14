package render_test

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func Fragments(styler *styler.Spec, size winsize.Winsize, frags []text.Fragment) string {
	var buffer strings.Builder

	lineSize := winsize.New(
		size.Rows,
		size.Cols,
	)

	for _, f := range frags {
		txt := helper.NewText(
			f.Text,
			f.Size(),
		)
		spec := styler.Apply(f.Spec, lineSize, txt)

		fragSize := text.FragmentMeasure(size.Cols, f)
		lineSize.Cols = lineSize.Cols.Sub(fragSize)

		buffer.WriteString(spec)
	}

	return buffer.String()
}
