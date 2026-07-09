package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Frag struct {
	Text string
	Atom atom.Atom
	Spec spec.Spec
}

func NewFrag(text string) *Frag {
	return &Frag{
		Text: text,
		Atom: atom.None,
		Spec: spec.Empty(),
	}
}

func EmptyFrag() *Frag {
	return NewFrag("")
}

func FragFromRunes(runes []rune) *Frag {
	return NewFrag(string(runes))
}

func FragFromMeta(other *Frag) *Frag {
	return EmptyFrag().CopyMeta(other)
}

func (f *Frag) CopyMeta(other *Frag) *Frag {
	f.Atom = other.Atom
	f.Spec = other.Spec
	return f
}

func (f *Frag) AddAtom(styles ...atom.Atom) *Frag {
	newAtom := atom.Merge(styles...)
	f.Atom = atom.Merge(f.Atom, newAtom)
	return f
}

func (f *Frag) AddSpec(styles ...spec.Spec) *Frag {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *Frag) Size() winsize.Cols {
	return runes.Measure(f.Text)
}

func (f *Frag) Clone() *Frag {
	clone := FragFromMeta(f)
	clone.Text = f.Text
	return clone
}

func FragsMeasure(cols winsize.Cols, frags ...Frag) winsize.Cols {
	measure := winsize.Cols(0)
	for _, f := range frags {
		ctx := spec.LayoutContext{
			SizeCols: cols,
			TextSize: f.Size(),
		}
		measure += spec.Measure(f.Spec, ctx)
	}
	return measure
}

func IsZeroFrag(frag Frag) bool {
	return frag.Text == "" &&
		frag.Atom == atom.None &&
		frag.Spec.Kind() == spec.KindNone
}

func IsStructuralFrag(frag Frag) bool {
	hasStyles := frag.Atom != atom.None || frag.Spec.Kind() != spec.KindNone
	return frag.Text == "" && hasStyles
}
