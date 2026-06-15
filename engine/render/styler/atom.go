package styler

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

type AtomStyler func(string) string

func pa(k atom.Atom, s AtomStyler) dict.Pair[atom.Atom, AtomStyler] {
	return dict.NewPair(k, s)
}

var Atoms = dict.NewInmutableLinkedMap(
	pa(atom.Lower, func(text string) string {
		return strings.ToLower(text)
	}),
	pa(atom.Upper, func(text string) string {
		return strings.ToUpper(text)
	}),
	pa(atom.Bold, func(text string) string {
		return text
	}),
	pa(atom.Select, func(text string) string {
		return text
	}),
)

type Atom struct {
	table *dict.LinkedMap[atom.Atom, AtomStyler]
}

func NewAtom() *Atom {
	instance := &Atom{}
	return instance.lazyInit()
}

func NewDefaultAtom() *Atom {
	return &Atom{
		table: Atoms.Clone(false),
	}
}

func (a *Atom) lazyInit() *Atom {
	if a.table != nil {
		return a
	}

	a.table = dict.NewLinkedMap[atom.Atom, AtomStyler]()
	return a
}

func (a *Atom) Push(pair ...dict.Pair[atom.Atom, AtomStyler]) *Atom {
	a.lazyInit()

	a.table.SetPairs(pair...)
	return a
}

func (a *Atom) Apply(text string, styles ...atom.Atom) string {
	a.lazyInit()

	merged := atom.Merge(styles...)

	for k, p := range a.table.All() {
		if merged.HasAny(k) {
			text = p(text)
		}
	}

	return text
}
