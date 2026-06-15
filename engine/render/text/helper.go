package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func MaxLineMeasure(cols winsize.Cols, lines ...Line) winsize.Cols {
	size := winsize.Cols(0)
	for _, l := range lines {
		measure := FragmentMeasure(cols, l.Text...)
		size = max(size, measure)
	}
	return size
}

func FragmentsFromString(text ...string) []Fragment {
	fragments := make([]Fragment, len(text))
	for i, v := range text {
		fragments[i] = *NewFragment(v)
	}
	return fragments
}

func LineJump() *Line {
	return &Line{
		Text: FragmentsFromString(""),
		Spec: spec.Cover(),
	}
}

func ApplyLineSpec(style spec.Spec, lines ...Line) []Line {
	for i := range lines {
		lines[i].SetSpec(style)
	}
	return lines
}

func LinesHasAtom(atm atom.Atom, lines ...Line) bool {
	for _, line := range lines {
		if FragsHasAtom(atm, line.Text...) {
			return true
		}
	}
	return false
}

func FragsHasAtom(atm atom.Atom, frags ...Fragment) bool {
	for _, v := range frags {
		if v.Atom.HasAny(atm) {
			return true
		}
	}
	return false
}

func EraseAtom(atm atom.Atom, lines ...Line) bool {
	for _, line := range lines {
		for _, v := range line.Text {
			v.Atom = atom.Erase(v.Atom, atm)
		}
	}
	return false
}

func CloneLines(lines ...Line) []Line {
	clones := make([]Line, len(lines))
	for i, v := range lines {
		clones[i] = *v.Clone()
	}
	return clones
}
