package sink

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

var specStylesTable = dict.NewInmutableLinkedMap(
	dict.P(spec.KindJustifyRight, sinkLinePaddingLeft),
	dict.P(spec.KindJustifyLeft, sinkLinePaddingRight),
	dict.P(spec.KindJustifyCenter, sinkLinePaddingCenter),
)

func sinkLinePaddingLeft(style spec.Kind, line *text.Line, _ winsize.Cols) *text.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec
	line.UnshiftFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingRight(style spec.Kind, line *text.Line, _ winsize.Cols) *text.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec
	line.PushFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingCenter(style spec.Kind, line *text.Line, cols winsize.Cols) *text.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec

	size := commons.Mapd(delSpec.Args()[spec.KeyJustifyCenterSize], cols)
	txt := delSpec.Args()[spec.KeyJustifyCenterText].Stringf()

	fragSize := text.FragmentMeasure(cols, line.Text...)

	available := size.Sub(fragSize)
	available = max(0, available)

	left := available / 2
	if left > 0 {
		paddLeft := spec.JustifyRight(left, txt)
		line.UnshiftFragments(
			*text.EmptyFragment().AddSpec(paddLeft),
		)
	}

	right := available.Sub(left)
	if right > 0 {
		paddRight := spec.JustifyLeft(right, txt)
		line.PushFragments(
			*text.EmptyFragment().AddSpec(paddRight),
		)
	}

	return line
}

func ApplySinks(line *text.Line, cols winsize.Cols) *text.Line {
	for k, t := range specStylesTable.All() {
		if !line.Spec.Kind().HasAny(k) {
			continue
		}
		line = t(k, line, cols)
	}
	return line
}
