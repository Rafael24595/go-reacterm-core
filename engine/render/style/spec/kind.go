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

type Descriptor struct {
	Kind Kind
	Name string
	Args []ArgKey
}

func init() {
	lookup = make(map[Kind]Descriptor, len(Registry))

	for _, d := range Registry {
		lookup[d.Kind] = d
	}
}

var lookup map[Kind]Descriptor

var Registry = [...]Descriptor{
	{
		Kind: KindJustifyRight,
		Args: []ArgKey{
			KeyJustifyRightSize,
			KeyJustifyRightText,
		},
	},
	{
		Kind: KindJustifyLeft,
		Name: "JustifyLeft",
		Args: []ArgKey{
			KeyJustifyLeftSize,
			KeyJustifyLeftText,
		},
	},
	{
		Kind: KindJustifyCenter,
		Name: "JustifyCenter",
		Args: []ArgKey{
			KeyJustifyCenterSize,
			KeyJustifyCenterText,
		},
	},
	{
		Kind: KindExtendLeft,
		Name: "ExtendLeft",
		Args: []ArgKey{
			KeyExtendLeftSize,
			KeyExtendLeftText,
		},
	},
	{
		Kind: KindExtendRight,
		Name: "ExtendRight",
		Args: []ArgKey{
			KeyExtendRightSize,
			KeyExtendRightText,
		},
	},
	{
		Kind: KindTruncateLeft,
		Name: "TruncateLeft",
		Args: []ArgKey{
			KeyTruncateLeftSize,
		},
	},
	{
		Kind: KindTruncateRight,
		Name: "TruncateRight",
		Args: []ArgKey{
			KeyTruncateRightSize,
		},
	},
	{
		Kind: KindFill,
		Name: "Fill",
		Args: []ArgKey{
			KeyFillSize,
		},
	},
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
