package atom

type Atom uint8

const (
	None Atom = 0
	Bold Atom = 1 << iota
	Upper
	Lower
	Select
	Focus
	Wrap
	Break
)

func (s Atom) Uint8() uint8 {
	return uint8(s)
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

func (s Atom) HasAny(styles ...Atom) bool {
	for _, style := range styles {
		if s&style != 0 {
			return true
		}
	}
	return false
}

func (s Atom) HasNone(styles ...Atom) bool {
	return !s.HasAny(styles...)
}
