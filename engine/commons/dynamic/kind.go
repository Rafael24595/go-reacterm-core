package dynamic

type kind uint8

const (
	// Nil represents the absence of a value.
	Nil    kind = 0
	// Bool represents a boolean value (true or false).
	Bool   kind = 1
	// String represents a string value.
	String kind = 2

	// Int represents an integer value.
	Int   kind = 3
	// Int8 represents an 8-bit integer value.
	Int8  kind = 4
	// Int16 represents a 16-bit integer value.
	Int16 kind = 5
	// Int32 represents a 32-bit integer value.
	Int32 kind = 6
	// Int64 represents a 64-bit integer value.
	Int64 kind = 7

	// Uint represents an unsigned integer value.
	Uint   kind = 8
	// Uint8 represents an 8-bit unsigned integer value.
	Uint8  kind = 9
	// Uint16 represents a 16-bit unsigned integer value.
	Uint16 kind = 10
	// Uint32 represents a 32-bit unsigned integer value.
	Uint32 kind = 11
	// Uint64 represents a 64-bit unsigned integer value.
	Uint64 kind = 12

	// Float32 represents a 32-bit floating-point value.
	Float32 kind = 13
	// Float64 represents a 64-bit floating-point value.
	Float64 kind = 14

	// Fallback represents an unknown or unsupported value type.
	Fallback kind = 15
)

// Uint8 casts the kind tag into a uint8 for hashing purposes.
func (k kind) Uint8() uint8 {
	return uint8(k)
}
