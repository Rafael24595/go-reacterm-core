package spec

// TODO: Make Lookup and Registry immutable.
type Descriptor struct {
	Kind Kind
	Name string
	Args []ArgKey
}

func init() {
	Lookup = make(map[Kind]Descriptor, len(Registry))

	for _, d := range Registry {
		Lookup[d.Kind] = d
	}
}

var Lookup map[Kind]Descriptor

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
