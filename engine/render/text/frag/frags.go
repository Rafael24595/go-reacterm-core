package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func Empty() Frag {
	return New("", atom.None, spec.Empty())
}

func TextAtom(text string, atoms ...atom.Atom) Frag {
	return New(text, atom.Merge(atoms...), spec.Empty())
}

func TextSpec(text string, specs ...spec.Spec) Frag {
	return New(text, atom.None, spec.Merge(specs...))
}

func FromString(runes string) Frag {
	return New(runes, atom.None, spec.Empty())
}

func FromRunes(runes []rune) Frag {
	return FromString(string(runes))
}

func FromAtom(atom atom.Atom) Frag {
	return New("", atom, spec.Empty())
}

func FromSpec(spec spec.Spec) Frag {
	return New("", atom.None, spec)
}

func FromMeta(atom atom.Atom, spec spec.Spec) Frag {
	return New("", atom, spec)
}
