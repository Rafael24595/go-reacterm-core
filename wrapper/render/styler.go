package wrapper_render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"

	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
)

var Atoms = []styler.AtomRule{
	{Atom: atom.Bold, Fn: toBold},
	{Atom: atom.Dim, Fn: toDim},
	{Atom: atom.Select, Fn: toSelect},
}

func toBold(text string) string {
	if text == "" {
		return text
	}
	return wrapper_ansi.Bold + text + wrapper_ansi.NormalWeight
}

func toDim(text string) string {
	if text == "" {
		return text
	}
	return wrapper_ansi.Dim + text + wrapper_ansi.NormalWeight
}

func toSelect(text string) string {
	if text == "" {
		return text
	}
	return wrapper_ansi.Reverse + text + wrapper_ansi.NoReverse
}
