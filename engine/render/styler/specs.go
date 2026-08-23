package styler

import (
	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

var specs = []SpecRule{
	{spec.KindFill, specFill},
	{spec.KindTruncateLeft, specTruncateLeft},
	{spec.KindTruncateRight, specTruncateRight},
	{spec.KindJustifyRight, specJustifyRight},
	{spec.KindJustifyLeft, specJustifyLeft},
	{spec.KindJustifyCenter, specJustifyCenter},
	{spec.KindExtendLeft, specExtendLeft},
	{spec.KindExtendRight, specExtendRight},
}

func deduplicateSpecs(specs []SpecRule) []SpecRule {
	cache := make(map[spec.Kind]int)
	rules := make([]SpecRule, 0, len(specs))

	for _, v := range specs {
		if ci, ok := cache[v.Kind]; ok {
			rules[ci] = v
			continue
		}

		cache[v.Kind] = len(rules)
		rules = append(rules, v)
	}

	return rules
}

func specFill(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return fill(spec, cols, text), true
}

func specTruncateLeft(spec spec.Spec, _ winsize.Cols, text format.Text) (string, bool) {
	return truncateLeft(spec, text), false
}

func specTruncateRight(spec spec.Spec, _ winsize.Cols, text format.Text) (string, bool) {
	return truncateRight(spec, text), false
}

func specJustifyRight(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return justifyRight(spec, cols, text), false
}

func specJustifyLeft(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return justifyLeft(spec, cols, text), false
}

func specJustifyCenter(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return justifyCenter(spec, cols, text), false
}

func specExtendLeft(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return extendLeft(spec, cols, text), false
}

func specExtendRight(spec spec.Spec, cols winsize.Cols, text format.Text) (string, bool) {
	return extendRight(spec, cols, text), false
}
