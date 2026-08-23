package styler

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type SpecStyler func(spec.Spec, winsize.Cols, format.Text) (string, bool)

type SpecRule struct {
	Kind spec.Kind
	Fn   SpecStyler
}

type Spec struct {
	table []SpecRule
}

func NewSpec() *Spec {
	instance := &Spec{}
	return instance.lazyInit()
}

func NewDefaultSpec() *Spec {
	return &Spec{
		table: deduplicateSpecs(specs),
	}
}

func (s *Spec) lazyInit() *Spec {
	if s.table != nil {
		return s
	}

	s.table = make([]SpecRule, 0)
	return s
}

func (s *Spec) Push(rules ...SpecRule) *Spec {
	s.lazyInit()

	s.table = deduplicateSpecs(
		append(s.table, rules...),
	)

	return s
}

func (s *Spec) Apply(style spec.Spec, size winsize.Winsize, text format.Text) string {
	s.lazyInit()

	kind := style.Kind()
	for _, r := range s.table {
		if !kind.HasAny(r.Kind) {
			continue
		}

		textData, exit := r.Fn(style, size.Cols, text)
		if exit {
			return textData
		}

		textSize := spec.MeasureOf(r.Kind, style, spec.LayoutContext{
			SizeCols: size.Cols,
			TextSize: text.Size,
		})

		text = format.NewText(
			textData,
			textSize,
		)
	}

	return text.Data
}

func fill(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	if text.Data == "" {
		text = format.TextFromString(marker.DefaultPaddingText)
	}

	size := dynamic.MapOr(args[spec.KeyFillSize], cols)
	size = min(cols, size)

	return format.PatternRight(size, text)
}

func truncateLeft(style spec.Spec, text format.Text) string {
	if text.Data == "" {
		return text.Data
	}

	args := style.Args()

	size := dynamic.MapOr[winsize.Cols](args[spec.KeyTruncateLeftSize], 0)
	size = max(1, size)

	ellipsis := format.NewEllipsis(
		args[spec.KeyTruncateEllipsisText].Text(),
		marker.DefaultElipsisSize,
	)

	return format.TruncateLeft(size, text, ellipsis)
}

func truncateRight(style spec.Spec, text format.Text) string {
	if text.Data == "" {
		return text.Data
	}

	args := style.Args()

	size := dynamic.MapOr[winsize.Cols](args[spec.KeyTruncateRightSize], 0)
	size = max(1, size)

	ellipsis := format.NewEllipsis(
		args[spec.KeyTruncateEllipsisText].Text(),
		marker.DefaultElipsisSize,
	)

	return format.TruncateRight(size, text, ellipsis)
}

func justifyRight(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	size := dynamic.MapOr(args[spec.KeyJustifyRightSize], cols)
	size = min(cols, size)

	filler := args[spec.KeyJustifyRightText].
		StringOr(marker.DefaultPaddingText)

	return format.JustifyRight(size, text, filler)
}

func justifyLeft(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	size := dynamic.MapOr(args[spec.KeyJustifyLeftSize], cols)
	size = min(cols, size)

	filler := args[spec.KeyJustifyLeftText].
		StringOr(marker.DefaultPaddingText)

	return format.JustifyLeft(size, text, filler)
}

func justifyCenter(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	size := dynamic.MapOr(args[spec.KeyJustifyCenterSize], cols)
	size = min(cols, size)

	filler := args[spec.KeyJustifyCenterText].
		StringOr(marker.DefaultPaddingText)

	return format.JustifyCenter(size, text, filler)
}

func extendLeft(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	size := dynamic.MapOr[winsize.Cols](args[spec.KeyExtendLeftSize], 0)
	filler := args[spec.KeyExtendLeftText].Text()

	if filler == "" {
		filler = text.Data
		text = format.EmptyText()
	}

	return format.ExtendLeft(min(cols, size), text, filler)
}

func extendRight(style spec.Spec, cols winsize.Cols, text format.Text) string {
	args := style.Args()

	size := dynamic.MapOr[winsize.Cols](args[spec.KeyExtendRightSize], 0)
	filler := args[spec.KeyExtendRightText].Text()

	if filler == "" {
		filler = text.Data
		text = format.EmptyText()
	}

	return format.ExtendRight(min(cols, size), text, filler)
}
