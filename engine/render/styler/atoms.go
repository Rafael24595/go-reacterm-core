package styler

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

var atoms = []AtomRule{
	{atom.Lower, toLower},
	{atom.Upper, toUpper},
	{atom.Bold, toDefault},
	{atom.Dim, toDefault},
	{atom.Select, toDefault},
	{atom.Focus, toDefault},
	{atom.Wrap, toDefault},
	{atom.Break, toDefault},
}

func deduplicateAtoms(specs []AtomRule) []AtomRule {
	cache := make(map[atom.Atom]int)
	rules := make([]AtomRule, 0, len(specs))

	for _, v := range specs {
		if ci, ok := cache[v.Atom]; ok {
			rules[ci] = v
			continue
		}

		cache[v.Atom] = len(rules)
		rules = append(rules, v)
	}

	return rules
}
