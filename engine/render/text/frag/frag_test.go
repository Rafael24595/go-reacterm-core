package frag

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func TestNew(t *testing.T) {
	spc := spec.AlignCenter()

	frg := New("hello", atom.Bold, spc)

	assert.Equal(t, "hello", frg.text)
	assert.Equal(t, atom.Bold, frg.atom)
	assert.DeepEqual(t, spc, frg.spec)
}

func TestMeasure(t *testing.T) {
	frg := FromString("hello")

	assert.Equal(t, 5, frg.Measure())
}

func TestClone(t *testing.T) {
	frg := New("hello", atom.Bold, spec.Empty())

	clone := frg.Clone()

	assert.Equal(t, frg.text, clone.text)
	assert.Equal(t, frg.atom, clone.atom)
	assert.DeepEqual(t, frg.spec, clone.spec)
	assert.Equal(t, frg.Hash(), clone.Hash())
}

func TestMeasure_Empty(t *testing.T) {
	assert.Equal(t, 0, Measure(80))
}

func TestMeasure_AddsMeasures(t *testing.T) {
	frags := []Frag{
		FromString("ab"),
		FromString("cde"),
	}

	assert.Equal(t, 5, Measure(80, frags...))
}

func TestIsZero(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		assert.True(t,
			IsZero(Empty()),
		)
	})

	t.Run("text", func(t *testing.T) {
		assert.False(t,
			IsZero(FromString("x")),
		)
	})

	t.Run("atom", func(t *testing.T) {
		assert.False(t,
			IsZero(FromAtom(atom.Bold)),
		)
	})

	t.Run("spec", func(t *testing.T) {
		assert.False(t,
			IsZero(FromSpec(spec.JustifyCenter(1))),
		)
	})
}

func TestIsStructural(t *testing.T) {
	t.Run("false when zero", func(t *testing.T) {
		assert.False(t,
			IsStructural(Empty()),
		)
	})

	t.Run("true because atom", func(t *testing.T) {
		assert.True(t,
			IsStructural(FromAtom(atom.Bold)),
		)
	})

	t.Run("true because spec", func(t *testing.T) {
		assert.True(t,
			IsStructural(FromSpec(spec.JustifyCenter(1))),
		)
	})

	t.Run("false because text", func(t *testing.T) {
		assert.False(t,
			IsStructural(New("x", atom.Bold, spec.JustifyCenter(1))),
		)
	})
}

func TestHash_LazyEvaluation(t *testing.T) {
	frg := New("golang", atom.Bold, spec.Empty())

	assert.False(t, frg.hashed)
	assert.Equal(t, 0, frg.hash)

	h1 := frg.Hash()
	assert.True(t, frg.hashed)
	assert.NotEqual(t, 0, h1)

	h2 := frg.Hash()
	assert.Equal(t, h1, h2)
}

func TestHashClone_LazyState(t *testing.T) {
	t.Run("Clone unhashed spec", func(t *testing.T) {
		original := New("zig", atom.Select, spec.Empty())
		clone := original.Clone()

		assert.False(t, original.hashed)
		assert.False(t, clone.hashed)

		assert.Equal(t, original.Hash(), clone.Hash())
	})

	t.Run("Clone already hashed spec", func(t *testing.T) {
		original := New("zig", atom.Select, spec.Empty())
		_ = original.Hash()

		clone := original.Clone()

		assert.True(t, clone.hashed)
		assert.Equal(t, original.hash, clone.hash)
		assert.Equal(t, original.Hash(), clone.Hash())
	})
}

func TestHash_Deterministic(t *testing.T) {
	frg1 := New("hello", atom.Bold, spec.Empty())
	frg2 := New("hello", atom.Bold, spec.Empty())

	assert.Equal(t, frg1.Hash(), frg2.Hash())
}

func TestHash_ChangesWhenContentChanges(t *testing.T) {
	base := New("hello", atom.Bold, spec.Empty())

	frg1 := New("world", atom.Bold, spec.Empty())
	assert.NotEqual(t, base.Hash(), frg1.Hash())

	frg2 := New("hello", atom.None, spec.Empty())
	assert.NotEqual(t, base.Hash(), frg2.Hash())

	frg3 := New("hello", atom.Bold, spec.JustifyCenter(1))
	assert.NotEqual(t, base.Hash(), frg3.Hash())
}

func BenchmarkNew(b *testing.B) {
	scp := spec.Empty()

	b.ReportAllocs()

	for b.Loop() {
		_ = New("hello world", atom.Bold, scp)
	}
}

func BenchmarkSize(b *testing.B) {
	cases := []struct {
		name string
		text string
	}{
		{"ascii10", "0123456789"},
		{"ascii100", strings.Repeat("a", 100)},
		{"unicode100", strings.Repeat("界", 100)},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			f := New(tc.text, atom.None, spec.Empty())

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = f.Measure()
			}
		})
	}
}

func BenchmarkCalcHash(b *testing.B) {
	scp := spec.Empty()

	b.ReportAllocs()

	for b.Loop() {
		_ = calcHash(
			hash.New(),
			"hello world",
			atom.Bold,
			scp,
		).Sum64()
	}
}
