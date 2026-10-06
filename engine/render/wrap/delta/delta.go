package delta

import "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"

// Delta represents a differential layout modification, tracking text fragments
// and word boundary offsets along with edge alignment flags.
type Delta struct {
	frags     []layout.Frag
	bounds    []uint32
	leftEdge  bool
	rightEdge bool
}

// New constructs and initializes a new, empty Delta instance.
func New() Delta {
	return Delta{
		frags:  make([]layout.Frag, 0),
		bounds: make([]uint32, 0),
	}
}
