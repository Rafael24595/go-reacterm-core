package spec

// Kind defines a bitmask representing composite visual transformation capabilities.
type Kind uint64

const (
	KindNone Kind = 0

	// Alignment kinds:

	// KindJustifyLeft aligns the text to the left.
	KindJustifyLeft Kind = 1 << (iota - 1)
	// KindJustifyRight aligns the text to the right.
	KindJustifyRight
	// KindJustifyCenter centers the text.
	KindJustifyCenter

	// Extension kinds:

	// KindExtendLeft extends the text to the left.
	KindExtendLeft
	// KindExtendRight extends the text to the right.
	KindExtendRight

	// Truncation kinds:

	// KindTruncateLeft truncates the text from the left.
	KindTruncateLeft
	// KindTruncateRight truncates the text from the right.
	KindTruncateRight

	// Fill kinds:

	// KindFill fills the available space with the specified text.
	KindFill
)

// Uint64 returns the primitive uint64 scalar representation of the Kind flag.
func (s Kind) Uint64() uint64 {
	return uint64(s)
}

// HasAny reports whether any of the specified target Kind flags are set.
func (s Kind) HasAny(styles ...Kind) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

// HasAll reports whether all of the specified target Kind flags are set.
func (s Kind) HasAll(styles ...Kind) bool {
	for _, style := range styles {
		if s&style == 0 {
			return false
		}
	}
	return true
}

// HasNone reports whether none of the specified target Kind flags are set.
func (s Kind) HasNone(styles ...Kind) bool {
	return !s.HasAny(styles...)
}

// ArgKey identifies specific key parameters inside a Spec's argument store.
type ArgKey uint8

const (
	// KeyJustifyLeftSize specifies the target size when left justifying.
	KeyJustifyLeftSize ArgKey = iota
	// KeyJustifyLeftText specifies the filler text when left justifying.
	KeyJustifyLeftText
	
	// KeyJustifyRightSize specifies the target size when right justifying.
	KeyJustifyRightSize
	// KeyJustifyRightText specifies the filler text when right justifying.
	KeyJustifyRightText

	// KeyJustifyCenterSize specifies the target size when centering.
	KeyJustifyCenterSize
	// KeyJustifyCenterText specifies the filler text when centering.
	KeyJustifyCenterText

	// KeyExtendLeftSize specifies the extension size to the left.
	KeyExtendLeftSize
	// KeyExtendLeftText specifies the filler text for left extension.
	KeyExtendLeftText

	// KeyExtendRightSize specifies the extension size to the right.
	KeyExtendRightSize
	// KeyExtendRightText specifies the filler text for right extension.
	KeyExtendRightText

	// KeyTruncateLeftSize specifies the maximum size after left truncation.
	KeyTruncateLeftSize
	// KeyTruncateRightSize specifies the maximum size after right truncation.
	KeyTruncateRightSize
	// KeyTruncateEllipsisText specifies the ellipsis text used for truncation.
	KeyTruncateEllipsisText

	// KeyFillSize specifies the target fill size.
	KeyFillSize
)

// Uint8 returns the primitive uint8 scalar representation of the ArgKey.
func (s ArgKey) Uint8() uint8 {
	return uint8(s)
}
