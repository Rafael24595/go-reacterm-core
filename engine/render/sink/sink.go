package sink

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type sinkFn func(style spec.Kind, lne line.Line, cols winsize.Cols) line.Line

type sinkRule struct {
	kind spec.Kind
	fn   sinkFn
}

var sinkipeline = [...]sinkRule{
	{spec.KindJustifyRight, sinkLinePaddingLeft},
	{spec.KindJustifyLeft, sinkLinePaddingRight},
	{spec.KindJustifyCenter, sinkLinePaddingCenter},
}

func sinkLinePaddingLeft(style spec.Kind, lne line.Line, _ winsize.Cols) line.Line {
	resSpec, delSpec := spec.Erase(lne.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return lne
	}

	return line.NewBuilder().
		WithLine(lne).
		SetSpec(resSpec).
		UnshiftFrags(frag.FromSpec(delSpec)).
		Line()
}

func sinkLinePaddingRight(style spec.Kind, lne line.Line, _ winsize.Cols) line.Line {
	resSpec, delSpec := spec.Erase(lne.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return lne
	}

	return line.NewBuilder().
		WithLine(lne).
		SetSpec(resSpec).
		PushFrags(frag.FromSpec(delSpec)).
		Line()
}

func sinkLinePaddingCenter(style spec.Kind, lne line.Line, cols winsize.Cols) line.Line {
	resSpec, delSpec := spec.Erase(lne.Spec, style)
	if delSpec.Kind() == spec.KindNone {
		return lne
	}

	builder := line.NewBuilder().
		WithLine(lne).
		SetSpec(resSpec)

	size := dynamic.MapOr(delSpec.Args()[spec.KeyJustifyCenterSize], cols)
	txt := delSpec.Args()[spec.KeyJustifyCenterText].Text()

	fragSize := frag.Measure(cols, builder.Text...)

	available := size.Sub(fragSize)
	available = max(0, available)

	left := available / 2
	if left > 0 {
		paddLeft := spec.JustifyRight(left, txt)
		builder.UnshiftFrags(
			frag.FromSpec(paddLeft),
		)
	}

	right := available.Sub(left)
	if right > 0 {
		paddRight := spec.JustifyLeft(right, txt)
		builder.PushFrags(
			frag.FromSpec(paddRight),
		)
	}

	return builder.Line()
}

func ApplySinks(line line.Line, cols winsize.Cols) line.Line {
	for _, t := range sinkipeline {
		if !line.Spec.Kind().HasAny(t.kind) {
			continue
		}
		line = t.fn(t.kind, line, cols)
	}
	return line
}
