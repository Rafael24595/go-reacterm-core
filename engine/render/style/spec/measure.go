package spec

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type LayoutContext struct {
	SizeCols winsize.Cols
	TextSize winsize.Cols
}

var measureTable = dict.NewInmutableLinkedMap(
	dict.P(KindFill, func(spep Spec, ctx LayoutContext) winsize.Cols {
		return commons.Mapd(spep.args[KeyFillSize], ctx.SizeCols)
	}),
	dict.P(KindTruncateLeft, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyTruncateLeftSize], ctx.TextSize)
		return min(ctx.TextSize, arg)
	}),
	dict.P(KindTruncateRight, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyTruncateRightSize], ctx.TextSize)
		return min(ctx.TextSize, arg)
	}),
	dict.P(KindJustifyCenter, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyJustifyCenterSize], ctx.SizeCols)
		return min(ctx.SizeCols, arg)
	}),
	dict.P(KindJustifyRight, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyJustifyRightSize], ctx.TextSize)
		return max(ctx.TextSize, arg)
	}),
	dict.P(KindJustifyLeft, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyJustifyLeftSize], ctx.TextSize)
		return max(ctx.TextSize, arg)
	}),
	dict.P(KindExtendLeft, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyExtendLeftSize], ctx.TextSize)
		return max(ctx.TextSize, arg)
	}),
	dict.P(KindExtendRight, func(spec Spec, ctx LayoutContext) winsize.Cols {
		arg := commons.Mapd(spec.args[KeyExtendRightSize], ctx.TextSize)
		return max(ctx.TextSize, arg)
	}),
)

func Measure(spec Spec, ctx LayoutContext) winsize.Cols {
	for k, c := range measureTable.All() {
		if spec.kind.HasAny(k) {
			ctx.TextSize = c(spec, ctx)
		}
	}
	return ctx.TextSize
}

func MeasureOf(kind Kind, spec Spec, ctx LayoutContext) winsize.Cols {
	if c, ok := measureTable.Get(kind); ok {
		return c(spec, ctx)
	}
	return ctx.TextSize
}
