package spec

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

const WinSize = winsize.Cols(0)

func Empty() Spec {
	return New(KindNone, args{})
}

func AlignLeft() Spec {
	return JustifyLeft(WinSize)
}

func AlignRight() Spec {
	return JustifyRight(WinSize)
}

func AlignCenter() Spec {
	return JustifyCenter(WinSize)
}

func JustifyLeft(size winsize.Cols, filler ...string) Spec {
	return specFromTextAndSize(
		KindJustifyLeft,
		KeyJustifyLeftSize,
		KeyJustifyLeftText,
		size,
		filler...,
	)
}

func JustifyRight(size winsize.Cols, filler ...string) Spec {
	return specFromTextAndSize(
		KindJustifyRight,
		KeyJustifyRightSize,
		KeyJustifyRightText,
		size,
		filler...,
	)
}

func JustifyCenter(size winsize.Cols, filler ...string) Spec {
	return specFromTextAndSize(
		KindJustifyCenter,
		KeyJustifyCenterSize,
		KeyJustifyCenterText,
		size,
		filler...,
	)
}

func ExtendLeft(size winsize.Cols, filler ...string) Spec {
	return specFromTextAndSize(
		KindExtendLeft,
		KeyExtendLeftSize,
		KeyExtendLeftText,
		size,
		filler...,
	)
}

func ExtendRight(size winsize.Cols, filler ...string) Spec {
	return specFromTextAndSize(
		KindExtendRight,
		KeyExtendRightSize,
		KeyExtendRightText,
		size,
		filler...,
	)
}

func TruncateLeft(size winsize.Cols, ellipsis ...string) Spec {
	spec := specFromSize(
		KindTruncateLeft,
		KeyTruncateLeftSize,
		size,
	)

	if len(ellipsis) > 0 {
		spec.args.Set(
			KeyTruncateEllipsisText,
			dynamic.From(ellipsis[0]),
		)
	}

	return spec
}

func TruncateRight(size winsize.Cols, ellipsis ...string) Spec {
	spec := specFromSize(
		KindTruncateRight,
		KeyTruncateRightSize,
		size,
	)

	if len(ellipsis) > 0 {
		spec.args.Set(
			KeyTruncateEllipsisText,
			dynamic.From(ellipsis[0]),
		)
	}

	return spec
}

func Cover() Spec {
	return Fill(WinSize)
}

func Fill(size winsize.Cols) Spec {
	return specFromSize(
		KindFill,
		KeyFillSize,
		size,
	)
}

func specFromSize(
	kind Kind,
	sizeKey ArgKey,
	size winsize.Cols,
) Spec {
	args := args{}

	if size != WinSize {
		args.Set(
			sizeKey, dynamic.From(size),
		)
	}

	return New(kind, args)
}

func specFromTextAndSize(
	kind Kind,
	sizeKey,
	textKey ArgKey,
	size winsize.Cols,
	text ...string,
) Spec {
	args := args{}

	if size != WinSize {
		args.Set(
			sizeKey, dynamic.From(size),
		)
	}

	if len(text) > 0 {
		args.Set(
			textKey,
			dynamic.From(strings.Join(text, "")),
		)
	}

	return New(kind, args)
}
