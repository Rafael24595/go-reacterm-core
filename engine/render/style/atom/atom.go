package atom

import "strings"

// Atom represents a bitmask combination of style, text, and structural flags.
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

// Uint16 returns the underlying uint16 representation of the Atom bitmask.
func (s Atom) Uint16() uint16 {
	return uint16(s)
}

// Merge combines multiple Atom bitmasks into a single composite Atom using bitwise OR.
func Merge(styles ...Atom) Atom {
	var merged Atom
	for _, style := range styles {
		merged |= style
	}
	return merged
}

// Erase removes the specified style flags from the target Atom bitmask using bitwise AND NOT.
func Erase(target, styles Atom) Atom {
	target &= ^styles
	return target
}

// HasAny returns true if the Atom contains at least one of the provided style flags.
func (a Atom) HasAny(styles ...Atom) bool {
	for _, style := range styles {
		if a&style != 0 {
			return true
		}
	}
	return false
}

// HasAll returns true if the Atom contains ALL of the provided style flags.
func (a Atom) HasAll(styles ...Atom) bool {
	for _, style := range styles {
		if a&style != style {
			return false
		}
	}
	return true
}

// HasNone returns true if the Atom contains NONE of the provided style flags.
func (s Atom) HasNone(styles ...Atom) bool {
	return !s.HasAny(styles...)
}

// String returns a human-readable representation of active flags for debugging.
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
