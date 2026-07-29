package wrap

type Delta struct {
	Frags  []wordFrag
	Bounds []uint32

	LeftEdge  bool
	RightEdge bool
}

func NewDelta() Delta {
	return Delta{
		Frags:  make([]wordFrag, 0),
		Bounds: make([]uint32, 0),
	}
}

func (f Delta) Size() uint32 {
	return uint32(len(f.Frags))
}

func (f *Delta) AddFrag(frag wordFrag) {
	f.Frags = append(f.Frags, frag)
}

func (f *Delta) BoundAtEnd() {
	f.Bounds = append(f.Bounds, f.Size())
}

func (f *Delta) Merge(other Delta) {
	if len(other.Frags) == 0 {
		return
	}

	if len(f.Frags) == 0 {
		f.Frags = append(f.Frags, other.Frags...)
		f.Bounds = append(f.Bounds, other.Bounds...)

		f.LeftEdge = other.LeftEdge
		f.RightEdge = other.RightEdge

		return
	}

	offset := f.Size()
	if f.RightEdge != other.LeftEdge {
		f.Bounds = append(f.Bounds, offset)
	}

	f.Frags = append(f.Frags, other.Frags...)

	for _, b := range other.Bounds {
		f.Bounds = append(f.Bounds, b+offset)
	}

	f.RightEdge = other.RightEdge
}
