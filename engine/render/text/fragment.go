package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

type Fragment struct {
	Text string
	Atom atom.Atom
	Spec spec.Spec
}

func NewFragment(text string) *Fragment {
	return &Fragment{
		Text: text,
		Atom: atom.None,
		Spec: spec.Empty(),
	}
}

func EmptyFragment() *Fragment {
	return NewFragment("")
}

func FragmentFromRunes(runes []rune) *Fragment {
	return NewFragment(string(runes))
}

func FragmentFromMeta(other *Fragment) *Fragment {
	return EmptyFragment().CopyMeta(other)
}

func (f *Fragment) CopyMeta(other *Fragment) *Fragment {
	f.Atom = other.Atom
	f.Spec = other.Spec
	return f
}

func (f *Fragment) AddAtom(styles ...atom.Atom) *Fragment {
	newAtom := atom.Merge(styles...)
	f.Atom = atom.Merge(f.Atom, newAtom)
	return f
}

func (f *Fragment) CutAtom(styles atom.Atom) *Fragment {
	f.Atom = atom.Erase(f.Atom, styles)
	return f
}

func (f *Fragment) AddSpec(styles ...spec.Spec) *Fragment {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *Fragment) CutSpec(styles spec.Kind) *Fragment {
	f.Spec, _ = spec.Erase(f.Spec, styles)
	return f
}

func (f *Fragment) Size() winsize.Cols {
	return runes.Measure(f.Text)
}

func (f *Fragment) Clone() *Fragment {
	clone := FragmentFromMeta(f)
	clone.Text = f.Text
	return clone
}

func FragmentMeasure(cols winsize.Cols, frags ...Fragment) winsize.Cols {
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

func IsZeroFragment(frag Fragment) bool {
	return frag.Text == "" &&
		frag.Atom == atom.None &&
		frag.Spec.Kind() == spec.KindNone
}

func IsStructuralFragment(frag Fragment) bool {
	hasStyles := frag.Atom != atom.None || frag.Spec.Kind() != spec.KindNone
	return frag.Text == "" && hasStyles
}
