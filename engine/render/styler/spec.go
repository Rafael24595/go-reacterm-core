package styler

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

type SpecStyler func(style.Spec, winsize.Cols, helper.Text) (string, bool)

func ps(k style.SpecKind, s SpecStyler) dict.Pair[style.SpecKind, SpecStyler] {
	return dict.NewPair(k, s)
}

var Specs = dict.NewInmutableLinkedMap(
	ps(style.SpcKindFill, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return fill(spec, cols, text), true
	}),
	ps(style.SpcKindTrimLeft, func(spec style.Spec, _ winsize.Cols, text helper.Text) (string, bool) {
		return trimLeft(spec, text), false
	}),
	ps(style.SpcKindTrimRight, func(spec style.Spec, _ winsize.Cols, text helper.Text) (string, bool) {
		return trimRight(spec, text), false
	}),
	ps(style.SpcKindPaddingCenter, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return paddingCenter(spec, cols, text), false
	}),
	ps(style.SpcKindPaddingLeft, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return paddingLeft(spec, cols, text), false
	}),
	ps(style.SpcKindPaddingRight, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return paddingRight(spec, cols, text), false
	}),
	ps(style.SpcKindRepeatLeft, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return repeatLeft(spec, cols, text), false
	}),
	ps(style.SpcKindRepeatRight, func(spec style.Spec, cols winsize.Cols, text helper.Text) (string, bool) {
		return repeatRight(spec, cols, text), false
	}),
)

type Spec struct {
	table *dict.LinkedMap[style.SpecKind, SpecStyler]
}

func NewSpec() *Spec {
	instance := &Spec{}
	return instance.lazyInit()
}

func NewDefaultSpec() *Spec {
	return &Spec{
		table: Specs.Clone(false),
	}
}

func (s *Spec) lazyInit() *Spec {
	if s.table != nil {
		return s
	}

	s.table = dict.NewLinkedMap[style.SpecKind, SpecStyler]()
	return s
}

func (s *Spec) Push(pair ...dict.Pair[style.SpecKind, SpecStyler]) *Spec {
	s.lazyInit()

	s.table.SetPairs(pair...)
	return s
}

func (s *Spec) Apply(spec style.Spec, size winsize.Winsize, text helper.Text) string {
	s.lazyInit()

	kind := spec.Kind()
	for k, p := range s.table.All() {
		if !kind.HasAny(k) {
			continue
		}

		textData, exit := p(spec, size.Cols, text)
		if exit {
			return textData
		}

		textSize := style.SpecMeasureOf(k, spec, style.LayoutContext{
			Cols:     size.Cols,
			TextSize: text.Size,
		})

		text = helper.NewText(
			textData,
			textSize,
		)
	}

	return text.Data
}

func fill(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	if text.Data == "" {
		text = helper.TextFromString(marker.DefaultPaddingText)
	}

	size := commons.Mapd(args[style.KeyFillSize], cols)
	size = min(cols, size)

	return helper.FillRight(size, text)
}

func trimLeft(spec style.Spec, text helper.Text) string {
	if text.Data == "" {
		return text.Data
	}

	args := spec.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyTrimLeftSize], 0)
	size = max(1, size)

	ellipsis := helper.NewEllipsis(
		args[style.KeyTrimEllipsisText].Stringf(),
		marker.DefaultElipsisSize,
	)

	return helper.TrimLeft(size, text, ellipsis)
}

func trimRight(spec style.Spec, text helper.Text) string {
	if text.Data == "" {
		return text.Data
	}

	args := spec.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyTrimRightSize], 0)
	size = max(1, size)

	ellipsis := helper.NewEllipsis(
		args[style.KeyTrimEllipsisText].Stringf(),
		marker.DefaultElipsisSize,
	)

	return helper.TrimRight(size, text, ellipsis)
}

func paddingCenter(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	size := commons.Mapd(args[style.KeyPaddingCenterSize], cols)
	size = min(cols, size)

	filler := args[style.KeyPaddingCenterText].
		Stringd(marker.DefaultPaddingText)

	return helper.Center(size, text, filler)
}

func paddingLeft(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	size := commons.Mapd(args[style.KeyPaddingLeftSize], cols)
	size = min(cols, size)

	filler := args[style.KeyPaddingLeftText].
		Stringd(marker.DefaultPaddingText)

	return helper.Left(size, text, filler)
}

func paddingRight(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	size := commons.Mapd(args[style.KeyPaddingRightSize], cols)
	size = min(cols, size)

	filler := args[style.KeyPaddingRightText].
		Stringd(marker.DefaultPaddingText)

	return helper.Right(size, text, filler)
}

func repeatLeft(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyRepeatLeftSize], 0)
	filler := args[style.KeyRepeatLeftText].Stringf()

	if filler == "" {
		filler = text.Data
		text = helper.EmptyText()
	}

	return helper.RepeatLeft(min(cols, size), text, filler)
}

func repeatRight(spec style.Spec, cols winsize.Cols, text helper.Text) string {
	args := spec.Args()

	size := commons.Mapd[winsize.Cols](args[style.KeyRepeatRightSize], 0)
	filler := args[style.KeyRepeatRightText].Stringf()

	if filler == "" {
		filler = text.Data
		text = helper.EmptyText()
	}

	return helper.RepeatRight(min(cols, size), text, filler)
}
