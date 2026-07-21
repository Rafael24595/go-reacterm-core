package stack

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestVStack_UnitBasicSuite(t *testing.T) {
	unit := VStackFromUnits()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestVStack_ToUnit_Anemic(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		PushWithOpts(mock.ToUnit()).
		ToUnit()

	assert.True(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, mock.Name, stack.Name)
}

func TestVStack_ToUnit_NotAnemic_MultipleElements(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Name: "mock_unit_001",
	}
	mock2 := &drawable_test.MockUnit{
		Name: "mock_unit_002",
	}

	stack := NewVStack().
		PushWithOpts(mock1.ToUnit()).
		PushWithOpts(mock2.ToUnit()).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}

func TestVStack_ToUnit_NotAnemic_LayerWithChunk(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		PushWithOpts(
			mock.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}

func TestVStack_ToUnit_NotAnemic_LayerStatic(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		PushWithOpts(
			mock.ToUnit(),
			layer.Static[winsize.Rows](),
		).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}

func TestVStack_ToUnit_NotAnemic_WithRenderer(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Name: "mock_unit",
	}

	stack := NewVStack().
		SetRenderer(defaultRenderer).
		PushWithOpts(mock.ToUnit()).
		ToUnit()

	assert.False(t, stack.Tags.Has(AnemicStack))
	assert.Equal(t, NameVStack, stack.Name)
}

func TestVStack_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock1 := &drawable_test.MockUnit{}

	unit := NewVStack().Push(
		mock1.ToUnit(),
	)

	unit.boot()

	assert.Panic(t, func() {
		m2 := &drawable_test.MockUnit{}
		unit.Push(m2.ToUnit())
	})
}

func TestVStack_Boot(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{}
	mock2 := &drawable_test.MockUnit{}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.boot()
	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.GreaterThan(t, 0, mock1.BootCalled)
	assert.GreaterThan(t, 0, mock2.BootCalled)
}

func TestVStack_Shift_Order(t *testing.T) {
	stack := &VStackUnit{}

	count := uint(0)

	mock1 := &drawable_test.MockUnit{Status: false}
	mock2 := &drawable_test.MockUnit{Status: false}

	unit1 := mock1.ToUnit()
	unit2 := mock2.ToUnit()

	unit1.Drawable.Draw = func(_ winsize.Winsize) ([]line.Line, bool) {
		mock1.DrawCalls = count
		count++
		return make([]line.Line, 0), false
	}

	unit2.Drawable.Draw = func(_ winsize.Winsize) ([]line.Line, bool) {
		mock2.DrawCalls = count
		count++
		return make([]line.Line, 0), false
	}

	stack.Push(unit1)
	stack.Push(unit2)

	stack.boot()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 0, mock1.DrawCalls)
	assert.Equal(t, 1, mock2.DrawCalls)
}

func TestVStack_Unshift_Order(t *testing.T) {
	stack := &VStackUnit{}

	count := uint(0)

	mock1 := &drawable_test.MockUnit{Status: false}
	mock2 := &drawable_test.MockUnit{Status: false}

	unit1 := mock1.ToUnit()
	unit2 := mock2.ToUnit()

	unit1.Drawable.Draw = func(_ winsize.Winsize) ([]line.Line, bool) {
		mock1.DrawCalls = count
		count++
		return make([]line.Line, 0), false
	}

	unit2.Drawable.Draw = func(_ winsize.Winsize) ([]line.Line, bool) {
		mock2.DrawCalls = count
		count++
		return make([]line.Line, 0), false
	}

	stack.Push(unit1)
	stack.Unshift(unit2)

	stack.boot()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 1, mock1.DrawCalls)
	assert.Equal(t, 0, mock2.DrawCalls)
}

func TestVStack_Draw_BreaksOnTrue(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{Status: true}
	mock2 := &drawable_test.MockUnit{Status: false}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.boot()

	_, global := stack.draw(winsize.Winsize{})

	assert.True(t, global)
	assert.Equal(t, 0, mock2.DrawCalls)
}

func TestVStack_DisablesLayer(t *testing.T) {
	stack := &VStackUnit{}

	mock := &drawable_test.MockUnit{Status: false}

	stack.Push(mock.ToUnit())

	stack.boot()

	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})
	stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Equal(t, 1, mock.DrawCalls)
}

func TestVStack_BufferConcat(t *testing.T) {
	stack := &VStackUnit{}

	line1 := line.FromString("go")
	line2 := line.FromString("lang")

	mock1 := &drawable_test.MockUnit{
		Lines:  []line.Line{line1},
		Status: false,
	}
	mock2 := &drawable_test.MockUnit{
		Lines:  []line.Line{line2},
		Status: false,
	}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
	)

	stack.boot()

	buffer, _ := stack.draw(winsize.Winsize{
		Rows: 10,
		Cols: 10,
	})

	assert.Size(t, 2, buffer)
	assert.Equal(
		t, "golang", text_test.LineToString(buffer[0])+text_test.LineToString(buffer[1]),
	)
}

func TestVStack_ShortCircuitStopsPropagation(t *testing.T) {
	stack := &VStackUnit{}

	mock1 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 1),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 2),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 1),
	}

	stack.Push(
		mock1.ToUnit(),
		mock2.ToUnit(),
		mock3.ToUnit(),
	)

	stack.boot()

	stack.draw(winsize.Winsize{
		Rows: 3,
		Cols: 10,
	})

	assert.Equal(t, 1, mock1.DrawCalls)
	assert.Equal(t, 1, mock2.DrawCalls)
	assert.Equal(t, 0, mock3.DrawCalls)
}

func TestVStack_FixedChunk_PadsWhenChildIsSmaller(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
	}

	stack := NewVStack().
		PushWithOpts(
			mock.ToUnit(),
			layer.Fixed[winsize.Rows](15),
		).
		ToUnit()

	stack.Drawable.Boot()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 20, Cols: 10})

	assert.Size(t, 15, lines)
}

func TestVStack_FixedChunk_TruncatesWhenChildIsBigger(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: make([]line.Line, 15),
	}

	stack := NewVStack().
		PushWithOpts(
			mock.ToUnit(),
			layer.Fixed[winsize.Rows](20),
		).
		ToUnit()

	stack.Drawable.Boot()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 10, Cols: 10})

	assert.Size(t, 10, lines)
}

func TestVStack_DynamicChunk_FillsRemainingSpace(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 5),
	}

	stack := NewVStack().
		PushWithOpts(
			mock1.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		PushWithOpts(mock2.ToUnit()).
		PushWithOpts(mock3.ToUnit()).
		ToUnit()

	stack.Drawable.Boot()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 30, Cols: 10})

	assert.Size(t, 25, lines)
}

func TestVStack_FixedOverflow_ShouldNotExceedContainer(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
	}

	stack := NewVStack().
		PushWithOpts(
			mock1.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		PushWithOpts(
			mock2.ToUnit(),
			layer.Fixed[winsize.Rows](10),
		).
		ToUnit()

	stack.Drawable.Boot()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 15, Cols: 10})

	assert.Size(t, 15, lines)
}

func TestVStack_ExactFit_NoExtraNoMissing(t *testing.T) {
	mock1 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 5),
	}
	mock2 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 5),
	}
	mock3 := &drawable_test.MockUnit{
		Lines: make([]line.Line, 5),
	}

	stack := NewVStack().
		PushWithOpts(mock1.ToUnit()).
		PushWithOpts(mock2.ToUnit()).
		PushWithOpts(mock3.ToUnit()).
		ToUnit()

	stack.Drawable.Boot()

	lines, _ := stack.Drawable.Draw(winsize.Winsize{Rows: 15, Cols: 10})

	assert.Size(t, 15, lines)
}
