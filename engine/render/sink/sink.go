package sink

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

var specStylesTable = dict.NewInmutableLinkedMap(
	dict.P(spec.KindJustifyRight, sinkLinePaddingLeft),
	dict.P(spec.KindJustifyLeft, sinkLinePaddingRight),
	dict.P(spec.KindJustifyCenter, sinkLinePaddingCenter),
)

func sinkLinePaddingLeft(style spec.Kind, line *line.Line, _ winsize.Cols) *line.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec
	line.UnshiftFrags(
		*frag.Empty().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingRight(style spec.Kind, line *line.Line, _ winsize.Cols) *line.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec
	line.PushFrags(
		*frag.Empty().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingCenter(style spec.Kind, line *line.Line, cols winsize.Cols) *line.Line {
	resSpec, delSpec := spec.Erase(line.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return line
	}

	line.Spec = resSpec

	size := dynamic.MapOr(delSpec.Args()[spec.KeyJustifyCenterSize], cols)
	txt := delSpec.Args()[spec.KeyJustifyCenterText].Text()

	fragSize := frag.Measure(cols, line.Text...)

	available := size.Sub(fragSize)
	available = max(0, available)

	left := available / 2
	if left > 0 {
		paddLeft := spec.JustifyRight(left, txt)
		line.UnshiftFrags(
			*frag.Empty().AddSpec(paddLeft),
		)
	}

	right := available.Sub(left)
	if right > 0 {
		paddRight := spec.JustifyLeft(right, txt)
		line.PushFrags(
			*frag.Empty().AddSpec(paddRight),
		)
	}

	return line
}

func ApplySinks(line *line.Line, cols winsize.Cols) *line.Line {
	for k, t := range specStylesTable.All() {
		if !line.Spec.Kind().HasAny(k) {
			continue
		}
		line = t(k, line, cols)
	}
	return line
}
