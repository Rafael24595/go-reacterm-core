package spec

import "iter"

type Descriptor struct {
	kind Kind
	name string
	args []ArgKey
}

func (d Descriptor) Kind() Kind {
	return d.kind
}

func (d Descriptor) Name() string {
	return d.name
}

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

func Lookup(kind Kind) (Descriptor, bool) {
	desc, ok := lookup[kind]
	return desc, ok
}

func Registry() iter.Seq[*Descriptor] {
	return func(yield func(*Descriptor) bool) {
		for i := range registry {
			if !yield(&registry[i]) {
				return
			}
		}
	}
}
