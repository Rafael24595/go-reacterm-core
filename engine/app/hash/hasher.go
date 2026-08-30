package hash

import "unsafe"

// Hasher represents an immutable 64-bit FNV-1a hash state.
type Hasher uint64

// Hash represents the computed 64-bit hash output.
type Hash uint64

type unsigned interface {
	~uint16 | ~uint32 | ~uint64
}

const (
	// FNV-1a 64-bit offset basis and prime constants.
	offset64 = 14695981039346656037
	// FNV-1a 64-bit prime constant.
	prime64  = 1099511628211
)

// New initializes a new Hasher with the standard FNV-1a offset basis.
func New() Hasher {
	return Hasher(offset64)
}

func writeUnsigned[T unsigned](h Hasher, v T) Hasher {
	for range unsafe.Sizeof(v) {
		h ^= Hasher(v & 0xFF)
		h *= prime64
		v >>= 8
	}
	return h
}

// Uint8 hashes a single byte into the state.
func (h Hasher) Uint8(v uint8) Hasher {
	h ^= Hasher(v)
	h *= prime64
	return h
}

// Uint16 hashes a 16-bit unsigned integer using Little-Endian byte order.
func (h Hasher) Uint16(v uint16) Hasher {
	return writeUnsigned(h, v)
}

// Uint32 hashes a 32-bit unsigned integer using Little-Endian byte order.
func (h Hasher) Uint32(v uint32) Hasher {
	return writeUnsigned(h, v)
}

// Uint64 hashes a 64-bit unsigned integer using Little-Endian byte order.
func (h Hasher) Uint64(v uint64) Hasher {
	return writeUnsigned(h, v)
}

// Hash combines an existing Hash value into the current Hasher state.
func (h Hasher) Hash(v Hash) Hasher {
	h ^= Hasher(v)
	h *= prime64
	return h
}

// Bool hashes a boolean value (1 for true, 0 for false).
func (h Hasher) Bool(v bool) Hasher {
	if v {
		return h.Uint8(1)
	}
	return h.Uint8(0)
}

// String hashes a raw string byte-by-byte without allocations.
func (h Hasher) String(s string) Hasher {
	for i := 0; i < len(s); i++ {
		h ^= Hasher(s[i])
		h *= prime64
	}
	return h
}

// Sum64 returns the finalized 64-bit Hash.
func (h Hasher) Sum64() Hash {
	return Hash(h)
}
