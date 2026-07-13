package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestNew(t *testing.T) {
    args := args{}
    args.Set(KeyFillSize, dynamic.From(winsize.Cols(10)))

    spec := New(KindFill, args)

    assert.Equal(t, KindFill, spec.kind)

    value, ok := spec.args.TryGet(KeyFillSize)
    assert.True(t, ok)
    assert.Equal(t, winsize.Cols(10), dynamic.MapOr(value, winsize.Cols(0)))

    assert.NotEqual(t, uint64(0), spec.hash)
}

func TestEmpty(t *testing.T) {
    spec := Empty()

    assert.Equal(t, KindNone, spec.kind)
    assert.Equal(t, 0, len(spec.args.items))
    assert.NotEqual(t, uint64(0), spec.hash)
}

func TestMerge_Empty(t *testing.T) {
    spec := Merge()

    assert.Equal(t, KindNone, spec.kind)
    assert.Equal(t, 0, len(spec.args.items))
}

func TestMerge_Single(t *testing.T) {
    original := Fill(20)

    merged := Merge(original)

    assert.Equal(t, original.kind, merged.kind)
    assert.Equal(t, original.hash, merged.hash)
    assert.Equal(t, len(original.args.items), len(merged.args.items))
}

func TestMerge_Multiple(t *testing.T) {
    spec := Merge(
        Fill(20),
        JustifyRight(10, "."),
        TruncateLeft(5),
    )

    assert.Equal(
        t,
        KindFill|KindJustifyRight|KindTruncateLeft,
        spec.kind,
    )

    assert.Inside(t, KeyFillSize, spec.args.items)
    assert.Inside(t, KeyJustifyRightSize, spec.args.items)
    assert.Inside(t, KeyJustifyRightText, spec.args.items)
    assert.Inside(t, KeyTruncateLeftSize, spec.args.items)
}

func TestMerge_LastSpecWins(t *testing.T) {
    spec := Merge(
        Fill(10),
        Fill(30),
    )

    value := dynamic.MapOr(
        spec.args.Get(KeyFillSize),
        winsize.Cols(0),
    )

    assert.Equal(t, winsize.Cols(30), value)
}

func TestSpecHash_IsOrderIndependent(t *testing.T) {
    left := Merge(
        Fill(20),
        JustifyRight(10, "."),
    )

    right := Merge(
        JustifyRight(10, "."),
        Fill(20),
    )

    assert.Equal(t, left.hash, right.hash)
}

func TestErase_All(t *testing.T) {
    original := Merge(
        Fill(10),
        JustifyRight(20, "."),
    )

    modified, removed := Erase(original, original.kind)

    assert.Equal(t, KindNone, modified.kind)
    assert.Equal(t, original.kind, removed.kind)
}

func TestEraseSpec_DeleteExists(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	modified, removed := Erase(scp, KindJustifyRight)
	size := dynamic.MapOr[winsize.Cols](removed.args.Get(KeyJustifyRightSize), 0)

	assert.Equal(t, KindFill, modified.kind)
	assert.NotInside(t, KeyJustifyRightSize, modified.args.items)
	assert.Inside(t, KeyFillSize, modified.args.items)

	assert.Equal(t, KindJustifyRight, removed.kind)
	assert.Equal(t, 10, size)
	assert.NotInside(t, KeyFillSize, removed.args.items)
}

func TestEraseSpec_DeleteNonExists(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	modified, removed := Erase(scp, KindTruncateLeft)

	assert.Equal(t, scp.kind, modified.kind)
	assert.Equal(t, len(scp.args.items), len(modified.args.items))

	assert.Equal(t, KindNone, removed.kind)
	assert.Equal(t, 0, len(removed.args.items))
}

func TestEraseSpec_DeleteMultiple(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	toRemove := KindJustifyRight | KindFill | KindExtendRight
	modified, removed := Erase(scp, toRemove)

	assert.Equal(t, KindNone, modified.kind)
	assert.Equal(t, 0, len(modified.args.items))

	assert.Equal(t, KindJustifyRight|KindFill, removed.kind)
	assert.Equal(t, 3, len(removed.args.items))
}

func TestErase_RebuildsHash(t *testing.T) {
    original := Merge(
        Fill(10),
        JustifyRight(20, "."),
    )

    modified, _ := Erase(original, KindFill)

    expected := Merge(
        JustifyRight(20, "."),
    )

    assert.Equal(t, expected.hash, modified.hash)
}

func TestSpecHash_IsDeterministic(t *testing.T) {
    left := Merge(
        Fill(20),
        JustifyRight(10, " "),
    )

    right := Merge(
        Fill(20),
        JustifyRight(10, " "),
    )

    assert.Equal(t, left.hash, right.hash)
}

func TestSpecHash_ChangesWhenArgumentChanges(t *testing.T) {
    left := Fill(20)
    right := Fill(21)

    assert.NotEqual(t, left.hash, right.hash)
}

func TestSpecHash_ChangesWhenKindChanges(t *testing.T) {
    left := Fill(20)
    right := TruncateLeft(20)

    assert.NotEqual(t, left.hash, right.hash)
}

func BenchmarkMeasure(b *testing.B) {
	spec := AlignCenter()
	ctx := LayoutContext{
		SizeCols: 80,
		TextSize: 20,
	}

	b.ReportAllocs()

	for b.Loop() {
		Measure(spec, ctx)
	}
}
