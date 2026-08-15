package delta

import "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"

type Builder struct {
	Frags     []layout.Frag
	Bounds    []uint32
	LeftEdge  bool
	RightEdge bool
}

func NewBuilder() *Builder {
	return &Builder{
		Frags:  make([]layout.Frag, 0),
		Bounds: make([]uint32, 0),
	}
}

func (b Builder) Size() uint32 {
	return uint32(len(b.Frags))
}

func (b *Builder) AddFrag(frag layout.Frag) *Builder {
	b.Frags = append(b.Frags, frag)
	return b
}

func (b *Builder) BoundAtEnd() *Builder {
	b.Bounds = append(b.Bounds, b.Size())
	return b
}

func (b *Builder) SetLeftEdge(leftEdge bool) *Builder {
	b.LeftEdge = leftEdge
	return b
}

func (b *Builder) SetRightEdge(rightEdge bool) *Builder {
	b.RightEdge = rightEdge
	return b
}

func (b *Builder) WithDelta(delta Delta) *Builder {
	if len(delta.frags) == 0 {
		return b
	}

	if len(b.Frags) == 0 {
		b.Frags = append(b.Frags, delta.frags...)
		b.Bounds = append(b.Bounds, delta.bounds...)

		b.LeftEdge = delta.leftEdge
		b.RightEdge = delta.rightEdge

		return b
	}

	offset := b.Size()
	if b.RightEdge != delta.leftEdge {
		b.Bounds = append(b.Bounds, offset)
	}

	b.Frags = append(b.Frags, delta.frags...)

	for _, v := range delta.bounds {
		b.Bounds = append(b.Bounds, v+offset)
	}

	b.RightEdge = delta.rightEdge

	return b
}

func (b *Builder) ToDelta() Delta {
	return Delta{
		frags:     b.Frags,
		bounds:    b.Bounds,
		leftEdge:  b.LeftEdge,
		rightEdge: b.RightEdge,
	}
}
