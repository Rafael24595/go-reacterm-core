package atom

import "iter"

// Descriptor holds metadata associating an Atom bitmask flag with its textual name.
type Descriptor struct {
	atom Atom
	name string
}

// Atom returns the Atom flag associated with this descriptor.
func (d Descriptor) Atom() Atom {
	return d.atom
}

// Name returns the textual representation/name of the Atom flag.
func (d Descriptor) Name() string {
	return d.name
}

func init() {
	lookup = make(map[Atom]Descriptor, len(registry))

	for _, d := range registry {
		lookup[d.atom] = d
	}
}

var lookup map[Atom]Descriptor

var registry = [...]Descriptor{
	{
		atom: Bold,
		name: "Bold",
	},
	{
		atom: Dim,
		name: "Dim",
	},
	{
		atom: Upper,
		name: "Upper",
	},
	{
		atom: Lower,
		name: "Lower",
	},
	{
		atom: Select,
		name: "Select",
	},
	{
		atom: Focus,
		name: "Focus",
	},
	{
		atom: Wrap,
		name: "Wrap",
	},
	{
		atom: Break,
		name: "Break",
	},
}

// Lookup searches for a registered Descriptor corresponding to the given Atom flag.
func Lookup(atom Atom) (Descriptor, bool) {
	desc, ok := lookup[atom]
	return desc, ok
}

// Registry returns an iterator yielding all registered Atom descriptors in order.
func Registry() iter.Seq[*Descriptor] {
	return func(yield func(*Descriptor) bool) {
		for i := range registry {
			if !yield(&registry[i]) {
				return
			}
		}
	}
}
