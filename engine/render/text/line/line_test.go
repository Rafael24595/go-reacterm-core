package line

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/hash"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestNew(t *testing.T) {
	ord := uint16(42)
	spc := spec.AlignCenter()
	frs := frag.FromStrings("hello", "golang")

	lne := New(ord, spc, frs)

	assert.Equal(t, ord, lne.order)
	assert.DeepEqual(t, spc, lne.spec)
	assert.DeepEqual(t, frs, lne.text)
}

func TestAt(t *testing.T) {
	frg0 := frag.FromString("golang")
	frg1 := frag.FromString("zig")

	zero := frag.Frag{}

	lne := New(
		1, spec.Empty(), []frag.Frag{frg0, frg1},
	)

	tests := []struct {
		name     string
		index    uint
		wantFrag frag.Frag
		wantOK   bool
	}{
		{"Valid first index", 0, frg0, true},
		{"Valid last index", 1, frg1, true},
		{"Out of bounds equal to size", 2, zero, false},
		{"Out of bounds far away", 99, zero, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := lne.At(tt.index)

			assert.Equal(t, tt.wantOK, ok)
			assert.DeepEqual(t, tt.wantFrag, got)
		})
	}
}

func TestLine_AtOrZero(t *testing.T) {
	frg := frag.FromString("golang")

	lne := New(
		1, spec.Empty(), []frag.Frag{frg},
	)

	t.Run("Valid index returns frag", func(t *testing.T) {
		got := lne.AtOrZero(0)
		assert.DeepEqual(t, frg, got)
	})

	t.Run("Out of bounds returns zero value", func(t *testing.T) {
		got := lne.AtOrZero(10)
		assert.DeepEqual(t, frag.Frag{}, got)
	})
}

func TestLine_Slice(t *testing.T) {
	frg := frag.FromString("golang")

	lne := New(
		1, spec.Empty(), []frag.Frag{frg},
	)

	slice := lne.Slice()

	assert.Size(t, 1, slice)

	slice[0] = frag.Frag{}
	got, _ := lne.At(0)

	assert.NotDeepEqual(t, slice[0], got)
	assert.NotSame(t, &slice[0], &got)
}

func TestLine_All(t *testing.T) {
	frgs := frag.FromStrings(
		"golang", "zig", "rust",
	)

	lne := New(1, spec.Empty(), frgs)

	t.Run("Iterates all elements", func(t *testing.T) {
		var collected []frag.Frag
		for f := range lne.All() {
			collected = append(collected, f)
		}
		assert.Size(t, len(frgs), collected)
	})

	t.Run("Respects early exit (break)", func(t *testing.T) {
		count := 0
		for range lne.All() {
			count++
			break
		}
		assert.Equal(t, 1, count)
	})
}

func TestLine_Clone(t *testing.T) {
	lne := New(
		42,
		spec.AlignCenter(),
		frag.FromStrings("hello", "golang"),
	)

	clone := lne.Clone()

	assert.Equal(t, clone.order, lne.order)
	assert.DeepEqual(t, clone.spec, lne.spec)
	assert.DeepEqual(t, clone.text, lne.text)
	assert.Equal(t, lne.hash, clone.hash)
}

func TestHash_Deterministic(t *testing.T) {
	spc := spec.Empty()
	frgs := frag.FromStrings("hello")

	l1 := New(1, spc, frgs)
	l2 := New(1, spc, frgs)

	assert.Equal(t, l1.hash, l2.hash)
}

func TestHash_ChangesWhenContentChanges(t *testing.T) {
	spc := spec.Empty()
	frgs := frag.FromStrings("hello")

	base := New(1, spc, frgs)

	assert.NotEqual(t, base.Hash(), New(2, spc, frgs).hash)
	assert.NotEqual(t, base.Hash(), New(1, spc, frag.FromStrings("world")).hash)
}

func TestFragsMeasure(t *testing.T) {
	l := New(1, spec.Empty(), []frag.Frag{
		frag.FromString("hello"),
		frag.FromString("world"),
	})

	cols := winsize.Cols(80)
	measure := frag.Measure(cols, l.Slice()...)

	assert.Equal(t, 10, measure)
}

func BenchmarkNew(b *testing.B) {
	spc := spec.Empty()
	frgs := frag.FromStrings(
		"hello", " ", "world",
	)

	b.ReportAllocs()

	for b.Loop() {
		_ = New(1, spc, frgs)
	}
}

func BenchmarkCalcHash(b *testing.B) {
	spc := spec.Empty()
	frgs := []frag.Frag{frag.FromString("hello world")}

	b.ReportAllocs()

	for b.Loop() {
		_ = calcHash(
			hash.New(),
			1,
			spc,
			frgs,
		).Sum64()
	}
}
