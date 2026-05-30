package indexmenu

type Pointer uint8

const (
	pointerBold Pointer = iota
	pointerSelect
)

var pointers = []Pointer{
	pointerSelect,
	pointerBold,
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
