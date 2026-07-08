package hash

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNew_ReturnsOffset(t *testing.T) {
	h := New()

	assert.Equal(t, offset64, h.Sum64())
}

func TestHasher_Uint8_IsDeterministic(t *testing.T) {
	h1 := New().Uint8(42).Sum64()
	h2 := New().Uint8(42).Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_Uint16_IsDeterministic(t *testing.T) {
	h1 := New().Uint16(42).Sum64()
	h2 := New().Uint16(42).Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_Uint32_IsDeterministic(t *testing.T) {
	h1 := New().Uint32(42).Sum64()
	h2 := New().Uint32(42).Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_Uint64_IsDeterministic(t *testing.T) {
	h1 := New().Uint64(42).Sum64()
	h2 := New().Uint64(42).Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_Uint64_DifferentValuesProduceDifferentHashes(t *testing.T) {
	h1 := New().Uint64(1).Sum64()
	h2 := New().Uint64(2).Sum64()

	assert.NotEqual(t, h1, h2)
}

func TestHasher_Bool(t *testing.T) {
	h1 := New().Bool(true).Sum64()
	h2 := New().Bool(false).Sum64()

	assert.NotEqual(t, h1, h2)
}

func TestHasher_Bool_EqualsUint8(t *testing.T) {
	h1 := New().Bool(true).Sum64()
	h2 := New().Uint8(1).Sum64()

	assert.Equal(t, h1, h2)

	h1 = New().Bool(false).Sum64()
	h2 = New().Uint8(0).Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_String_Empty(t *testing.T) {
	h := New().String("").Sum64()

	assert.Equal(t, offset64, h)
}

func TestHasher_String_IsDeterministic(t *testing.T) {
	h1 := New().String("golang").Sum64()
	h2 := New().String("golang").Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasherString_IsOrderSensitive(t *testing.T) {
	h1 := New().String("go").Sum64()
	h2 := New().String("og").Sum64()

	assert.NotEqual(t, h1, h2)
}

func TestHasherString_IsCaseSensitive(t *testing.T) {
	h1 := New().String("Golang").Sum64()
	h2 := New().String("golang").Sum64()

	assert.NotEqual(t, h1, h2)
}

func TestHasher_Chaining(t *testing.T) {
	h1 := New().
		String("golang").
		Uint16(42).
		Bool(true).
		Sum64()

	h2 := New().
		String("golang").
		Uint16(42).
		Bool(true).
		Sum64()

	assert.Equal(t, h1, h2)
}

func TestHasher_OrderMatters(t *testing.T) {
	h1 := New().
		String("golang").
		Uint8(1).
		Sum64()

	h2 := New().
		Uint8(1).
		String("golang").
		Sum64()

	assert.NotEqual(t, h1, h2)
}

func TestHasher_IsImmutable(t *testing.T) {
	base := New()

	_ = base.Uint64(42)

	assert.Equal(t, uint64(offset64), base.Sum64())
}

func TestHasher_GoldenValue(t *testing.T) {
	got := New().
		String("golang").
		Uint16(80).
		Uint16(24).
		Bool(true).
		Sum64()

	assert.Equal(t, 0xdbf288cb32d28296, got)
}

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = New()
	}
}

func BenchmarkUint8(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.Uint8(42)
	}
}

func BenchmarkUint16(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.Uint16(42)
	}
}

func BenchmarkUint32(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.Uint32(42)
	}
}

func BenchmarkUint64(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.Uint64(42)
	}
}

func BenchmarkBool(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.Bool(true)
	}
}

func BenchmarkStringShort(b *testing.B) {
	h := New()

	for b.Loop() {
		_ = h.String("golang")
	}
}

func BenchmarkStringLongASCII(b *testing.B) {
	s := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 64)
	h := New()

	b.SetBytes(int64(len(s)))

	for b.Loop() {
		_ = h.String(s)
	}
}

func BenchmarkStringUnicode(b *testing.B) {
	s := strings.Repeat("áéí你好🙂", 64)
	h := New()

	b.SetBytes(int64(len(s)))

	for b.Loop() {
		_ = h.String(s)
	}
}

func BenchmarkFragmentLikeHash(b *testing.B) {
	for b.Loop() {
		_ = New().
			String("Hello world").
			Uint8(3).
			Uint64(0x123456789abcdef).
			Sum64()
	}
}

func BenchmarkLineHash(b *testing.B) {
	frags := []struct {
		text string
		atom uint8
		spec uint64
	}{
		{"Hello", 1, 10},
		{" ", 0, 0},
		{"world", 2, 15},
		{"!", 0, 0},
	}

	b.ReportAllocs()

	for b.Loop() {
		h := New()

		for _, f := range frags {
			h = h.
				String(f.text).
				Uint8(f.atom).
				Uint64(f.spec)
		}

		_ = h.Sum64()
	}
}
