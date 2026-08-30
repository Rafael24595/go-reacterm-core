package hash

import "unsafe"

type Hasher uint64

type Hash uint64

type unsigned interface {
	~uint16 | ~uint32 | ~uint64
}

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

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

func (h Hasher) Uint8(v uint8) Hasher {
	h ^= Hasher(v)
	h *= prime64
	return h
}

func (h Hasher) Uint16(v uint16) Hasher {
	return writeUnsigned(h, v)
}

func (h Hasher) Uint32(v uint32) Hasher {
	return writeUnsigned(h, v)
}

func (h Hasher) Uint64(v uint64) Hasher {
	return writeUnsigned(h, v)
}

func (h Hasher) Hash(v Hash) Hasher {
	h ^= Hasher(v)
	h *= prime64
	return h
}

func (h Hasher) Bool(v bool) Hasher {
	if v {
		return h.Uint8(1)
	}
	return h.Uint8(0)
}

func (h Hasher) String(s string) Hasher {
	for i := 0; i < len(s); i++ {
		h ^= Hasher(s[i])
		h *= prime64
	}
	return h
}

func (h Hasher) Sum64() Hash {
	return Hash(h)
}
