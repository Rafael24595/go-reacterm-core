package argument

type kind uint8

const (
	Nil    kind = 0
	Bool   kind = 1
	String kind = 2

	Int   kind = 3
	Int8  kind = 4
	Int16 kind = 5
	Int32 kind = 6
	Int64 kind = 7

	Uint   kind = 8
	Uint8  kind = 9
	Uint16 kind = 10
	Uint32 kind = 11
	Uint64 kind = 12

	Float32 kind = 13
	Float64 kind = 14

	Fallback kind = 15
)

func (k kind) Uint8() uint8 {
	return uint8(k)
}
