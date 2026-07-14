package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type BuilderLine struct {
	Order uint16
	Text  []frag.Frag
	Spec  spec.Spec
}

func (f *BuilderLine) SetOrder(order uint16) *BuilderLine {
	f.Order = order
	return f
}

func (f *BuilderLine) AddFrags(text ...frag.Frag) *BuilderLine {
	f.Text = append(f.Text, text...)
	return f
}

func (f *BuilderLine) AddBuilder(builder ...*frag.Builder) *BuilderLine {
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
