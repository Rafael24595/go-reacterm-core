package wrapper_render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/dict"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"

	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
)

func pa(k atom.Atom, s styler.AtomStyler) dict.Pair[atom.Atom, styler.AtomStyler] {
	return dict.NewPair(k, s)
}

var Atoms = dict.NewInmutableLinkedMap(
	pa(atom.Bold, func(text string) string {
		if text == "" {
			return text
		}
		return wrapper_ansi.Bold + text + wrapper_ansi.NormalWeight
	}),
	pa(atom.Select, func(text string) string {
		if text == "" {
			return text
		}
		return wrapper_ansi.Reverse + text + wrapper_ansi.NoReverse
	}),
)
