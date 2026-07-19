package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type Builder struct {
	Order uint16
	Spec  spec.Spec
	Text  []frag.Frag
}

func (b *Builder) SetOrder(order uint16) *Builder {
	b.Order = order
	return b
}

func (b *Builder) AddSpec(styles ...spec.Spec) *Builder {
	newSpec := spec.Merge(styles...)
	b.Spec = spec.Merge(b.Spec, newSpec)
	return b
}

func (b *Builder) UnshiftFrags(frags ...frag.Frag) *Builder {
	b.Text = append(frags, b.Text...)
	return b
}

func (b *Builder) PushFrags(frags ...frag.Frag) *Builder {
	b.Text = append(b.Text, frags...)
	return b
}

func (b *Builder) UnshiftBuilder(builder ...*frag.Builder) *Builder {
	frags := make([]frag.Frag, len(builder))
	for i := range builder {
		frags[i] = builder[i].Frag()
	}
	return b.PushFrags(frags...)
}

func (b *Builder) PushBuilder(builder ...*frag.Builder) *Builder {
	for _, f := range builder {
		b.Text = append(b.Text, f.Frag())
	}
	return b
}

func (b *Builder) Line() Line {
	return Line{
		Order: b.Order,
		Spec:  b.Spec,
		Text:  b.Text,
	}
}
