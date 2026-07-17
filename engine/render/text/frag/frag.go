package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Frag struct {
	text string
	atom atom.Atom
	spec spec.Spec
	hash uint64
}

func New(
	text string,
	atom atom.Atom,
	spec spec.Spec,
) Frag {
	hash := calcHash(
		hash.New(),
		text,
		atom,
		spec,
	)

	return Frag{
		text: text,
		atom: atom,
		spec: spec,
		hash: hash.Sum64(),
	}
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

func (f Frag) Text() string {
	return f.text
}

func (f Frag) Atom() atom.Atom {
	return f.atom
}

func (f Frag) Spec() spec.Spec {
	return f.spec
}

func (s Frag) Hash() uint64 {
	return s.hash
}

func (f Frag) Size() winsize.Cols {
	return runes.Measure(f.text)
}

func (f Frag) Clone() Frag {
	return New(
		f.text,
		f.atom,
		f.spec,
	)
}

func Measure(cols winsize.Cols, frags ...Frag) winsize.Cols {
	measure := winsize.Cols(0)
	for _, f := range frags {
		ctx := spec.LayoutContext{
			SizeCols: cols,
			TextSize: f.Size(),
		}
		measure += spec.Measure(f.spec, ctx)
	}
	return measure
}

func IsZero(frag Frag) bool {
	return frag.text == "" &&
		frag.atom == atom.None &&
		frag.spec.Kind() == spec.KindNone
}

func IsStructural(frag Frag) bool {
	hasStyles := frag.atom != atom.None || frag.spec.Kind() != spec.KindNone
	return frag.text == "" && hasStyles
}
