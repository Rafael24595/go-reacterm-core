package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func BreakWord(frags []text.Frag) []text.Frag {
	for i := range frags {
		frags[i].AddAtom(atom.Break)
	}
	return frags
}
