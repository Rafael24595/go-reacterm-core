package styler

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

type AtomStyler func(string) string

type AtomRule struct {
	Atom atom.Atom
	Fn   AtomStyler
}

func toLower(text string) string {
	return strings.ToLower(text)
}

func toUpper(text string) string {
	return strings.ToUpper(text)
}

func toDefault(text string) string {
	return text
}

type Atom struct {
	table []AtomRule
}

func NewAtom() *Atom {
	instance := &Atom{}
	return instance.lazyInit()
}

func NewDefaultAtom() *Atom {
	return &Atom{
		table: deduplicateAtoms(atoms),
	}
}

func (a *Atom) lazyInit() *Atom {
	if a.table != nil {
		return a
	}

	a.table = make([]AtomRule, 0)
	return a
}

func (a *Atom) Push(rules ...AtomRule) *Atom {
	a.lazyInit()

	a.table = deduplicateAtoms(
		append(a.table, rules...),
	)

	return a
}

func (a *Atom) Apply(text string, styles ...atom.Atom) string {
	a.lazyInit()

	merged := atom.Merge(styles...)

	for _, r := range a.table {
		if merged.HasAny(r.Atom) {
			text = r.Fn(text)
		}
	}

	return text
}
