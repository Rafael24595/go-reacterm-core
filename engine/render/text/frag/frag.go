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
		Text: text,
		Atom: atom,
		Spec: spec,
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

func (f Frag) Size() winsize.Cols {
	return runes.Measure(f.Text)
}

func (f Frag) Clone() Frag {
	return New(
		f.Text,
		f.Atom,
		f.Spec,
	)
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
