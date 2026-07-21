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

func NewBuilder(capacity ...int) *Builder {
	size := 0
	if len(capacity) > 0 {
		size = capacity[0]
	}

	return &Builder{
		Order: 0,
		Spec:  spec.Empty(),
		Text:  make([]frag.Frag, 0, size),
	}
}

func BuilderFromLine(lne Line) *Builder {
	return NewBuilder().
		WithLine(lne)
}

func (b *Builder) SetOrder(order uint16) *Builder {
	b.Order = order
	return b
}

func (b *Builder) SetSpec(styles ...spec.Spec) *Builder {
	b.Spec = spec.Merge(styles...)
	return b
}

func (b *Builder) AddSpec(styles ...spec.Spec) *Builder {
	newSpec := spec.Merge(styles...)
	b.Spec = spec.Merge(b.Spec, newSpec)
	return b
}

func (b *Builder) WithMeta(other Line) *Builder {
	b.Order = other.Order
	b.AddSpec(other.Spec)
	return b
}

func (b *Builder) WithLine(other Line) *Builder {
	b.Order = other.Order
	b.AddSpec(other.Spec)
	b.Text = append(b.Text, other.Text...)
	return b
}

func (b *Builder) UnshiftText(text ...string) *Builder {
	return b.UnshiftFrags(
		frag.FromStrings(text...)...,
	)
}

func (b *Builder) PushText(text ...string) *Builder {
	return b.PushFrags(
		frag.FromStrings(text...)...,
	)
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
	return b.UnshiftFrags(frags...)
}

func (b *Builder) PushBuilder(builder ...*frag.Builder) *Builder {
	for _, f := range builder {
		b.Text = append(b.Text, f.Frag())
	}
	return b
}

func (b *Builder) Line() Line {
	return New(
		b.Order,
		b.Spec,
		b.Text,
	)
}

func (b *Builder) LinePtr() *Line {
	lne := b.Line()
	return &lne
}
