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

func (b *Builder) AddRunes(runes []rune) *Builder {
	return b.AddText(
		string(runes),
	)
}

func (b *Builder) AddText(text string) *Builder {
	b.Text += text
	return b
}

func (b *Builder) WithMeta(other *Frag) *Builder {
	b.Atom = other.Atom
	b.Spec = other.Spec
	return b
}

func (b *Builder) WithFrag(frg Frag) *Builder {
	b.Text += frg.Text
	b.Atom = atom.Merge(b.Atom, frg.Atom)
	b.Spec = spec.Merge(b.Spec, frg.Spec)
	return b
}

func (b *Builder) AddAtom(styles ...atom.Atom) *Builder {
	newAtom := atom.Merge(styles...)
	b.Atom = atom.Merge(b.Atom, newAtom)
	return b
}

func (b *Builder) AddSpec(styles ...spec.Spec) *Builder {
	newSpec := spec.Merge(styles...)
	b.Spec = spec.Merge(b.Spec, newSpec)
	return b
}

func (b *Builder) Frag() Frag {
	return New(b.Text, b.Atom, b.Spec)
}
