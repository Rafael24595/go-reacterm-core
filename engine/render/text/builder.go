package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type BuilderFrag struct {
	Text string
	Atom atom.Atom
	Spec spec.Spec
}

func NewBuilderFrag() *BuilderFrag {
	return &BuilderFrag{
		Text: "",
		Atom: atom.None,
		Spec: spec.Empty(),
	}
}

func (f *BuilderFrag) AddRunes(runes []rune) *BuilderFrag {
	return f.AddText(
		string(runes),
	)
}

func (f *BuilderFrag) AddText(text string) *BuilderFrag {
	f.Text += text
	return f
}

func (f *BuilderFrag) WithMeta(other *Frag) *BuilderFrag {
	f.Atom = other.Atom
	f.Spec = other.Spec
	return f
}

func (f *BuilderFrag) AddAtom(styles ...atom.Atom) *BuilderFrag {
	newAtom := atom.Merge(styles...)
	f.Atom = atom.Merge(f.Atom, newAtom)
	return f
}

func (f *BuilderFrag) AddSpec(styles ...spec.Spec) *BuilderFrag {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *BuilderFrag) Frag() Frag {
	return Frag{
		Text: f.Text,
		Atom: f.Atom,
		Spec: f.Spec,
	}
}

type BuilderLine struct {
	Order uint16
	Text  []Frag
	Spec  spec.Spec
}

func (f *BuilderLine) SetOrder(order uint16) *BuilderLine {
	f.Order = order
	return f
}

func (f *BuilderLine) AddFrags(text ...Frag) *BuilderLine {
	f.Text = append(f.Text, text...)
	return f
}

func (f *BuilderLine) AddBuilder(builder ...*BuilderFrag) *BuilderLine {
	for _, b := range builder {
		f.Text = append(f.Text, b.Frag())
	}
	return f
}

func (f *BuilderLine) AddSpec(styles ...spec.Spec) *BuilderLine {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *BuilderLine) Frag() Line {
	return Line{
		Order: f.Order,
		Text:  f.Text,
		Spec:  f.Spec,
	}
}
