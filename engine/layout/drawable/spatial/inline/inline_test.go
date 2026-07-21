package inline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestInline_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := Wrap(mock.ToUnit())
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestInline_JoinsChildren(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: []line.Line{
			line.FromString("go"),
		},
	}
	mock2 := &drawable_test.MockUnit{
		Lines: []line.Line{
			line.FromString("lang"),
		},
	}

	unit := New(
		mock1.ToUnit(),
		mock2.ToUnit(),
	).ToUnit()

	unit.Drawable.Boot()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Size(t, 1, lines)
	assert.Equal(t, "golang", text_test.LineToString(lines[0]))
}

func TestInline_JoinsChildrenWithSeparator(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: []line.Line{
			line.FromString("golang"),
		},
	}
	mock2 := &drawable_test.MockUnit{
		Lines: []line.Line{
			line.FromString("ziglang"),
		},
	}

	unit := New(
		mock1.ToUnit(),
		mock2.ToUnit(),
	).Separator(" | ").ToUnit()

	unit.Drawable.Boot()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 16,
	})

	assert.Size(t, 1, lines)
	assert.Equal(t, "golang | ziglang", text_test.LineToString(lines[0]))
}

func TestInline_MultipleLines(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			line.FromString("go"),
			line.FromString("lang"),
		},
	}

	unit := New(
		mock.ToUnit(),
	).Separator(" | ").ToUnit()

	unit.Drawable.Boot()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 3,
		Cols: 9,
	})

	assert.Size(t, 1, lines)
	assert.Equal(t, "go | lang", text_test.LineToString(lines[0]))
}

func TestInline_Empty(t *testing.T) {
	unit := New().ToUnit()

	unit.Drawable.Boot()

	lines, _ := unit.Drawable.Draw(winsize.Winsize{})

	assert.Empty(t, lines)
}
