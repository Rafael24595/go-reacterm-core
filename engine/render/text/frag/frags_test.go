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

	assert.Equal(t, "hello", frg.Text())
	assert.Equal(t, atom.Bold, frg.Atom())
	assert.DeepEqual(t, spc, frg.Spec())
}

func TestSize(t *testing.T) {
	frg := FromString("hello")

	assert.Equal(t, 5, frg.Size())
}

func TestClone(t *testing.T) {
	frg := New("hello", atom.Bold, spec.Empty())

	clone := frg.Clone()

	assert.Equal(t, frg.Text(), clone.Text())
	assert.Equal(t, frg.Atom(), clone.Atom())
	assert.DeepEqual(t, frg.Spec(), clone.Spec())
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

func TestHash_Deterministic(t *testing.T) {
	frg1 := New("hello", atom.Bold, spec.Empty())
	frg2 := New("hello", atom.Bold, spec.Empty())

	assert.Equal(t, frg1.hash, frg2.hash)
}

func TestHash_ChangesWhenContentChanges(t *testing.T) {
	base := New("hello", atom.Bold, spec.Empty())

	assert.NotEqual(t, base.hash, New("world", atom.Bold, spec.Empty()).hash)
	assert.NotEqual(t, base.hash, New("hello", atom.None, spec.Empty()).hash)
	assert.NotEqual(t, base.hash, New("hello", atom.Bold, spec.JustifyCenter(1)).hash)
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
				_ = f.Size()
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
