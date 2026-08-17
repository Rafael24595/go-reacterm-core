package spec

type Kind uint64

const (
	KindNone Kind = 0

	KindJustifyLeft Kind = 1 << iota
	KindJustifyRight
	KindJustifyCenter

	KindExtendLeft
	KindExtendRight

	KindTruncateLeft
	KindTruncateRight

	KindFill
)

func (s Kind) Uint64() uint64 {
	return uint64(s)
}

func (s Kind) HasAny(styles ...Kind) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s Kind) HasNone(styles ...Kind) bool {
	return !s.HasAny(styles...)
}

type ArgKey uint8

const (
	KeyJustifyLeftSize ArgKey = iota
	KeyJustifyLeftText

	KeyJustifyRightSize
	KeyJustifyRightText

	KeyJustifyCenterSize
	KeyJustifyCenterText

	KeyExtendLeftSize
	KeyExtendLeftText

	KeyExtendRightSize
	KeyExtendRightText

	KeyTruncateLeftSize
	KeyTruncateRightSize
	KeyTruncateEllipsisText

	KeyFillSize
)

func (s ArgKey) Uint8() uint8 {
	return uint8(s)
}
