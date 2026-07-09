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
