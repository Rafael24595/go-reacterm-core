package atom

// TODO: Make Lookup and Registry immutable.
type Descriptor struct {
	Atom Atom
	Name string
}

func init() {
	Lookup = make(map[Atom]Descriptor, len(Registry))

	for _, d := range Registry {
		Lookup[d.Atom] = d
	}
}

var Lookup map[Atom]Descriptor

var Registry = [...]Descriptor{
	{
		Atom: Bold,
		Name: "Bold",
	},
	{
		Atom: Dim,
		Name: "Dim",
	},
	{
		Atom: Upper,
		Name: "Upper",
	},
	{
		Atom: Lower,
		Name: "Lower",
	},
	{
		Atom: Select,
		Name: "Select",
	},
	{
		Atom: Focus,
		Name: "Focus",
	},
	{
		Atom: Wrap,
		Name: "Wrap",
	},
	{
		Atom: Break,
		Name: "Break",
	},
}
