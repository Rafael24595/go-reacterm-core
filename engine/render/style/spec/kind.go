package spec

type KindDescriptor struct {
	Kind Kind
	Args []ArgKey
}

func init() {
	kindLookup = make(map[Kind]KindDescriptor, len(kindRegistry))

	for _, d := range kindRegistry {
		kindLookup[d.Kind] = d
	}
}

var kindLookup map[Kind]KindDescriptor

var kindRegistry = [...]KindDescriptor{
	{
		Kind: KindJustifyRight,
		Args: []ArgKey{
			KeyJustifyRightSize,
			KeyJustifyRightText,
		},
	},
	{
		Kind: KindJustifyLeft,
		Args: []ArgKey{
			KeyJustifyLeftSize,
			KeyJustifyLeftText,
		},
	},
	{
		Kind: KindJustifyCenter,
		Args: []ArgKey{
			KeyJustifyCenterSize,
			KeyJustifyCenterText,
		},
	},
	{
		Kind: KindExtendLeft,
		Args: []ArgKey{
			KeyExtendLeftSize,
			KeyExtendLeftText,
		},
	},
	{
		Kind: KindExtendRight,
		Args: []ArgKey{
			KeyExtendRightSize,
			KeyExtendRightText,
		},
	},
	{
		Kind: KindTruncateLeft,
		Args: []ArgKey{
			KeyTruncateLeftSize,
		},
	},
	{
		Kind: KindTruncateRight,
		Args: []ArgKey{
			KeyTruncateRightSize,
		},
	},
	{
		Kind: KindFill,
		Args: []ArgKey{
			KeyFillSize,
		},
	},
}

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

