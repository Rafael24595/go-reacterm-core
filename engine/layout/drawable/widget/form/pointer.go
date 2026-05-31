package form

type Pointer uint8

const (
	PointerPrompt Pointer = 1 << iota
	PointerGutter
)

func (m Pointer) HasAny(pointers ...Pointer) bool {
	for _, mod := range pointers {
		if m&mod != 0 {
			return true
		}
	}
	return false
}

func (m Pointer) HasNone(pointers ...Pointer) bool {
	return !m.HasAny(pointers...)
}

var pointers = []Pointer{
	PointerGutter,
	PointerPrompt,
	PointerPrompt | PointerGutter,
}

func FindPointer(cursor uint8) Pointer {
	if cursor >= uint8(len(pointers)) {
		return pointers[0]
	}
	return pointers[cursor]
}

func NextPointer(cursor uint8) uint8 {
	return (cursor + 1) % uint8(len(pointers))
}
