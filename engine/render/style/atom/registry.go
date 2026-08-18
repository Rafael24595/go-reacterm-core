package atom

import "iter"

type Descriptor struct {
	atom Atom
	name string
}

func (d Descriptor) Atom() Atom {
	return d.atom
}

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

func Lookup(atom Atom) (Descriptor, bool) {
	desc, ok := lookup[atom]
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
