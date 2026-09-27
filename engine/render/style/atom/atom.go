package atom

import "strings"

type Atom uint16

const (
	None Atom = 0

	// Text atoms:

	// Upper renders the element in uppercase.
	Upper Atom = 1 << (iota - 1)
	// Lower renders the element in lowercase.
	Lower

	// Style atoms:

	// Bold renders the element in bold.
	Bold
	// Dim renders the element with reduced intensity.
	Dim
	// Select renders the element as selected.
	Select

	// Structural atoms:

	// Focus indicates that the element is focused.
	Focus
	// Wrap enables text wrapping.
	Wrap
	// Break forces a text break.
	Break
)

func (s Atom) Uint16() uint16 {
	return uint16(s)
}

func Merge(styles ...Atom) Atom {
	var merged Atom
	for _, style := range styles {
		merged |= style
	}
	return merged
}

func Erase(target, styles Atom) Atom {
	target &= ^styles
	return target
}

func (a Atom) HasAny(styles ...Atom) bool {
	for _, style := range styles {
		if a&style != 0 {
			return true
		}
	}
	return false
}

func (a Atom) HasAll(styles ...Atom) bool {
	for _, style := range styles {
		if a&style != style {
			return false
		}
	}
	return true
}

func (s Atom) HasNone(styles ...Atom) bool {
	return !s.HasAny(styles...)
}

func (a Atom) String() string {
	if a == None {
		return "None"
	}

	var parts []string
	flags := []struct {
		flag Atom
		name string
	}{
		{Upper, "Upper"},
		{Lower, "Lower"},
		{Bold, "Bold"},
		{Dim, "Dim"},
		{Select, "Select"},
		{Focus, "Focus"},
		{Wrap, "Wrap"},
		{Break, "Break"},
	}

	for _, f := range flags {
		if a.HasAny(f.flag) {
			parts = append(parts, f.name)
		}
	}

	return strings.Join(parts, "|")
}
