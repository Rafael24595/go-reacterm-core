package spec

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestMeasure_Fill(t *testing.T) {
	s := Fill(80)

	ctx := LayoutContext{
		SizeCols: winsize.Cols(120),
		TextSize: winsize.Cols(10),
	}

	res := Measure(s, ctx)
	assert.Equal(t, 80, res)
}

func TestMeasure_Truncate(t *testing.T) {
	s := TruncateRight(15)

	ctx := LayoutContext{
		SizeCols: winsize.Cols(100),
		TextSize: winsize.Cols(30),
	}

	res := Measure(s, ctx)
	assert.Equal(t, 15, res)
}

func TestMeasure_Justify(t *testing.T) {
	s := JustifyRight(40, " ")

	ctx := LayoutContext{
		SizeCols: winsize.Cols(100),
		TextSize: winsize.Cols(20),
	}

	res := Measure(s, ctx)
	assert.Equal(t, 40, res)
}

func TestMeasureOf_SingleKind(t *testing.T) {
	s := Merge(
		Fill(50),
		JustifyRight(80, " "),
	)

	ctx := LayoutContext{
		SizeCols: winsize.Cols(100),
		TextSize: winsize.Cols(20),
	}

	resFill := MeasureOf(KindFill, s, ctx)
	assert.Equal(t, 50, resFill)

	resJustify := MeasureOf(KindJustifyRight, s, ctx)
	assert.Equal(t, 80, resJustify)
}
