package transformer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func BreakWord(frags []frag.Frag) []frag.Frag {
	for i := range frags {
		frags[i] = frag.NewBuilder().
			WithFrag(frags[i]).
			AddAtom(atom.Break).
			Frag()
	}
	return frags
}
