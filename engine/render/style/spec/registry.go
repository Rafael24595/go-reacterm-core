package spec

import "iter"

// Descriptor defines metadata associating a layout Kind flag with its expected argument keys.
type Descriptor struct {
	kind Kind
	name string
	args []ArgKey
}

// Kind returns the layout Kind bitmask flag.
func (d Descriptor) Kind() Kind {
	return d.kind
}

// Name returns the textual identifier of the layout transformation.
func (d Descriptor) Name() string {
	return d.name
}

// Args returns an iterator yielding all ArgKey dependencies for this Kind.
func (d Descriptor) Args() iter.Seq[ArgKey] {
	return func(yield func(ArgKey) bool) {
		for i := range d.args {
			if !yield(d.args[i]) {
				return
			}
		}
	}
}

func init() {
	lookup = make(map[Kind]Descriptor, len(registry))

	for _, d := range registry {
		lookup[d.kind] = d
	}
}

var lookup map[Kind]Descriptor

var registry = [...]Descriptor{
	{
		kind: KindJustifyRight,
		name: "JustifyRight",
		args: []ArgKey{
			KeyJustifyRightSize,
			KeyJustifyRightText,
		},
	},
	{
		kind: KindJustifyLeft,
		name: "JustifyLeft",
		args: []ArgKey{
			KeyJustifyLeftSize,
			KeyJustifyLeftText,
		},
	},
	{
		kind: KindJustifyCenter,
		name: "JustifyCenter",
		args: []ArgKey{
			KeyJustifyCenterSize,
			KeyJustifyCenterText,
		},
	},
	{
		kind: KindExtendLeft,
		name: "ExtendLeft",
		args: []ArgKey{
			KeyExtendLeftSize,
			KeyExtendLeftText,
		},
	},
	{
		kind: KindExtendRight,
		name: "ExtendRight",
		args: []ArgKey{
			KeyExtendRightSize,
			KeyExtendRightText,
		},
	},
	{
		kind: KindTruncateLeft,
		name: "TruncateLeft",
		args: []ArgKey{
			KeyTruncateLeftSize,
		},
	},
	{
		kind: KindTruncateRight,
		name: "TruncateRight",
		args: []ArgKey{
			KeyTruncateRightSize,
		},
	},
	{
		kind: KindFill,
		name: "Fill",
		args: []ArgKey{
			KeyFillSize,
		},
	},
}

// Lookup searches for a registered Descriptor corresponding to the given Kind.
func Lookup(kind Kind) (Descriptor, bool) {
	desc, ok := lookup[kind]
	return desc, ok
}

// Registry returns an iterator yielding all registered layout spec descriptors.
func Registry() iter.Seq[*Descriptor] {
	return func(yield func(*Descriptor) bool) {
		for i := range registry {
			if !yield(&registry[i]) {
				return
			}
		}
	}
}
