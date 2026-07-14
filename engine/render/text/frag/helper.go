package frag

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

func FromStrings(text ...string) []Frag {
	frags := make([]Frag, len(text))
	for i, v := range text {
		frags[i] = *New(v)
	}
	return frags
}

func HasAtom(atm atom.Atom, frags ...Frag) bool {
	for _, v := range frags {
		if v.Atom.HasAny(atm) {
			return true
		}
	}
	return false
}
