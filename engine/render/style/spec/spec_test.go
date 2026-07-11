package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/argument"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestEraseSpec_DeleteExists(t *testing.T) {
	scp := Merge(
		JustifyRight(10, " "),
		Fill(100),
	)

	modified, removed := Erase(scp, KindJustifyRight)
	size := argument.Mapd[winsize.Cols](removed.args.Get(KeyJustifyRightSize), 0)

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
