package delta

import "github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"

type Delta struct {
	frags     []layout.Frag
	bounds    []uint32
	leftEdge  bool
	rightEdge bool
}

func New() Delta {
	return Delta{
		frags:  make([]layout.Frag, 0),
		bounds: make([]uint32, 0),
	}
}
