package sink

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

var specStylesTable = dict.NewInmutableLinkedMap(
	dict.P(style.SpcKindPaddingLeft, sinkLinePaddingLeft),
	dict.P(style.SpcKindPaddingRight, sinkLinePaddingRight),
	dict.P(style.SpcKindPaddingCenter, sinkLinePaddingCenter),
)

func sinkLinePaddingLeft(spec style.SpecKind, line *text.Line, _ winsize.Cols) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec
	line.UnshiftFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingRight(spec style.SpecKind, line *text.Line, _ winsize.Cols) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec
	line.PushFragments(
		*text.EmptyFragment().AddSpec(delSpec),
	)

	return line
}

func sinkLinePaddingCenter(spec style.SpecKind, line *text.Line, cols winsize.Cols) *text.Line {
	resSpec, delSpec := style.EraseSpec(line.Spec, spec)
	if delSpec.Kind() == style.SpcKindNone {
		return line
	}

	line.Spec = resSpec

	size := commons.Mapd(delSpec.Args()[style.KeyPaddingCenterSize], cols)
	txt := delSpec.Args()[style.KeyPaddingCenterText].Stringf()

	fragSize := text.FragmentMeasure(cols, line.Text...)

	available := size.Sub(fragSize)
	available = max(0, available)

	left := available / 2
	if left > 0 {
		paddLeft := style.SpecPaddingLeft(left, txt)
		line.UnshiftFragments(
			*text.EmptyFragment().AddSpec(paddLeft),
		)
	}

	right := available.Sub(left)
	if right > 0 {
		paddRight := style.SpecPaddingRight(right, txt)
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
