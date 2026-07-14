package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Frag struct {
	Text string
	Atom atom.Atom
	Spec spec.Spec
	hash uint64
}

func New(text string) *Frag {
	return newFrag(
		text,
		atom.None,
		spec.Empty(),
	)
}

func newFrag(
	text string,
	atom atom.Atom,
	spec spec.Spec,
) *Frag {
	hash := calcHash(
		hash.New(),
		text,
		atom,
		spec,
	)

	return &Frag{
		Text: text,
		Atom: atom,
		Spec: spec,
		hash: hash.Sum64(),
	}
}

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

func calcHash(
	hasher hash.Hasher,
	text string,
	atom atom.Atom,
	spec spec.Spec,
) hash.Hasher {
	hasher = hasher.String(text)
	hasher = hasher.Uint8(atom.Uint8())
	hasher = hasher.Uint64(spec.Hash())
	return hasher
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
	clone := FromMeta(f)
	clone.Text = f.Text
	return clone
}

func Measure(cols winsize.Cols, frags ...Frag) winsize.Cols {
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

func IsZero(frag Frag) bool {
	return frag.Text == "" &&
		frag.Atom == atom.None &&
		frag.Spec.Kind() == spec.KindNone
}

func IsStructural(frag Frag) bool {
	hasStyles := frag.Atom != atom.None || frag.Spec.Kind() != spec.KindNone
	return frag.Text == "" && hasStyles
}
