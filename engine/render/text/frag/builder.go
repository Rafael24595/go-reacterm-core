package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Builder struct {
	Text string
	Atom atom.Atom
	Spec spec.Spec
}

func NewBuilder() *Builder {
	return &Builder{
		Text: "",
		Atom: atom.None,
		Spec: spec.Empty(),
	}
}

func (f *Builder) AddRunes(runes []rune) *Builder {
	return f.AddText(
		string(runes),
	)
}

func (f *Builder) AddText(text string) *Builder {
	f.Text += text
	return f
}

func (f *Builder) WithMeta(other *Frag) *Builder {
	f.Atom = other.Atom
	f.Spec = other.Spec
	return f
}

func (f *Builder) AddAtom(styles ...atom.Atom) *Builder {
	newAtom := atom.Merge(styles...)
	f.Atom = atom.Merge(f.Atom, newAtom)
	return f
}

func (f *Builder) AddSpec(styles ...spec.Spec) *Builder {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *Builder) Frag() Frag {
	return Frag{
		Text: f.Text,
		Atom: f.Atom,
		Spec: f.Spec,
	}
}
