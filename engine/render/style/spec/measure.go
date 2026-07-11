package spec

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type LayoutContext struct {
	SizeCols winsize.Cols
	TextSize winsize.Cols
}

type measureFn func(Spec, LayoutContext) winsize.Cols

type measureRule struct {
	kind Kind
	fn   measureFn
}

func init() {
	measureLookup = make(map[Kind]measureFn, len(measurePipeline))

	for _, r := range measurePipeline {
		measureLookup[r.kind] = r.fn
	}
}

var measureLookup map[Kind]measureFn

var measurePipeline = [...]measureRule{
	{KindFill, measureFill},
	{KindTruncateLeft, measureTruncateLeft},
	{KindTruncateRight, measureTruncateRight},
	{KindJustifyCenter, measureJustifyCenter},
	{KindJustifyRight, measureJustifyRight},
	{KindJustifyLeft, measureJustifyLeft},
	{KindExtendLeft, measureExtendLeft},
	{KindExtendRight, measureExtendRight},
}

func measureFill(spec Spec, ctx LayoutContext) winsize.Cols {
	return dynamic.MapOr(spec.args.Get(KeyFillSize), ctx.SizeCols)
}

func measureTruncateLeft(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyTruncateLeftSize), ctx.TextSize)
	return min(ctx.TextSize, arg)
}

func measureTruncateRight(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyTruncateRightSize), ctx.TextSize)
	return min(ctx.TextSize, arg)
}

func measureJustifyCenter(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyJustifyCenterSize), ctx.SizeCols)
	return min(ctx.SizeCols, arg)
}

func measureJustifyRight(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyJustifyRightSize), ctx.TextSize)
	return max(ctx.TextSize, arg)
}

func measureJustifyLeft(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyJustifyLeftSize), ctx.TextSize)
	return max(ctx.TextSize, arg)
}

func measureExtendLeft(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyExtendLeftSize), ctx.TextSize)
	return max(ctx.TextSize, arg)
}

func measureExtendRight(spec Spec, ctx LayoutContext) winsize.Cols {
	arg := dynamic.MapOr(spec.args.Get(KeyExtendRightSize), ctx.TextSize)
	return max(ctx.TextSize, arg)
}

func Measure(spec Spec, ctx LayoutContext) winsize.Cols {
	for i := range measurePipeline {
		rule := &measurePipeline[i]

		if spec.kind.HasAny(rule.kind) {
			ctx.TextSize = rule.fn(spec, ctx)
		}
	}

	return ctx.TextSize
}

func MeasureOf(kind Kind, spec Spec, ctx LayoutContext) winsize.Cols {
	if fn, ok := measureLookup[kind]; ok {
		return fn(spec, ctx)
	}
	return ctx.TextSize
}
