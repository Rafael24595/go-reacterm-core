package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func Empty() *Frag {
	return FromString("")
}

func FromString(runes string) *Frag {
	return newFrag(runes, atom.None, spec.Empty())
}

func FromRunes(runes []rune) *Frag {
	return FromString(string(runes))
}

func FromMeta(other *Frag) *Frag {
	return Empty().CopyMeta(other)
}
