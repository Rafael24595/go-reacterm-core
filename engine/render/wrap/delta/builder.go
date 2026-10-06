package delta

import "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"

// Builder construct and accumulates layout fragments and boundary state
// to create a Delta instance.
type Builder struct {
	// Frags holds the list of layout fragments accumulated in the builder.
	Frags     []layout.Frag
	// Bounds holds the list of word boundary offsets corresponding to the fragments.
	Bounds    []uint32
	// LeftEdge indicates whether the left edge boundary is active.
	LeftEdge  bool
	// RightEdge indicates whether the right edge boundary is active.
	RightEdge bool
}

// NewBuilder constructs and initializes a new, empty Builder instance.
func NewBuilder() *Builder {
	return &Builder{
		Frags:  make([]layout.Frag, 0),
		Bounds: make([]uint32, 0),
	}
}

// Size returns the total count of fragments currently stored in the builder.
func (b Builder) Size() uint32 {
	return uint32(len(b.Frags))
}

// AddFrag appends a layout fragment to the builder.
func (b *Builder) AddFrag(frag layout.Frag) *Builder {
	b.Frags = append(b.Frags, frag)
	return b
}

// BoundAtEnd records a word boundary offset at the current end of the fragment list.
func (b *Builder) BoundAtEnd() *Builder {
	b.Bounds = append(b.Bounds, b.Size())
	return b
}

// SetLeftEdge sets the flag indicating whether the left edge boundary is active.
func (b *Builder) SetLeftEdge(leftEdge bool) *Builder {
	b.LeftEdge = leftEdge
	return b
}

// SetRightEdge sets the flag indicating whether the right edge boundary is active.
func (b *Builder) SetRightEdge(rightEdge bool) *Builder {
	b.RightEdge = rightEdge
	return b
}

// WithDelta merges the contents and boundaries of an existing Delta into the builder,
// recalculating word boundary offsets and aligning edges appropriately.
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

// ToDelta constructs and returns a Delta value populated with the current builder state.
func (b *Builder) ToDelta() Delta {
	return Delta{
		frags:     b.Frags,
		bounds:    b.Bounds,
		leftEdge:  b.LeftEdge,
		rightEdge: b.RightEdge,
	}
}
