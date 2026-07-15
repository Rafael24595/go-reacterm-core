package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type Builder struct {
	Order uint16
	Text  []frag.Frag
	Spec  spec.Spec
}

func (f *Builder) SetOrder(order uint16) *Builder {
	f.Order = order
	return f
}

func (f *Builder) AddFrags(text ...frag.Frag) *Builder {
	f.Text = append(f.Text, text...)
	return f
}

func (f *Builder) AddBuilder(builder ...*frag.Builder) *Builder {
	for _, b := range builder {
		f.Text = append(f.Text, b.Frag())
	}
	return f
}

func (f *Builder) AddSpec(styles ...spec.Spec) *Builder {
	newSpec := spec.Merge(styles...)
	f.Spec = spec.Merge(f.Spec, newSpec)
	return f
}

func (f *Builder) Line() Line {
	return Line{
		Order: f.Order,
		Text:  f.Text,
		Spec:  f.Spec,
	}
}
