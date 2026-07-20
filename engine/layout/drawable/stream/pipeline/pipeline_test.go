package pipeline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func mockInitStep(s winsize.Winsize, d drawable.Unit) drawable.Unit {
	return d
}

func mockDrawStep(s winsize.Winsize, d drawable.Unit) ([]line.Line, bool) {
	return d.Drawable.Draw(s)
}

func mockDataStep(_ winsize.Winsize, _ drawable.Unit, l []line.Line, s bool) ([]line.Line, bool) {
	return l, s
}

func TestPipeline_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		SetDrawStep(mockDrawStep).
		ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestPipeline_ShouldPanicIfNewElementsAddedAfterInitialization(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		SetDrawStep(mockDrawStep)

	unit.ToUnit().Drawable.Boot()

	assert.Panic(t, func() {
		unit.PushBootSteps(mockInitStep)
	})

	assert.Panic(t, func() {
		unit.SetDrawStep(mockDrawStep)
	})

	assert.Panic(t, func() {
		unit.PushDataSteps(mockDataStep)
	})
}

func TestPipeline_ReturnBaseIfNils(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).
		ToUnit()

	assert.Equal(t, drawable_test.NameMockUnit, unit.Name)

	unit = New(mock.ToUnit()).
		SetDrawStep(mockDrawStep).
		ToUnit()

	assert.Equal(t, Name, unit.Name)
}

func TestPipeline_BootStepTransformation(t *testing.T) {
	mock1 := &drawable_test.MockUnit{}

	mock2 := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("base_01"),
			*line.New("base_02"),
		},
		Status: true,
	}

	unit := New(mock1.ToUnit()).
		PushBootSteps(func(_ winsize.Winsize, _ drawable.Unit) drawable.Unit {
			return mock2.ToUnit()
		}).
		ToUnit()

	unit.Drawable.Boot()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Size(t, 2, lines)
	assert.True(t, status)
	assert.Equal(
		t, text_test.LineToString(mock2.Lines[0]), text_test.LineToString(lines[0]),
	)
}

func TestPipeline_DrawStepTransformation(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("base_01"),
			*line.New("base_02"),
			*line.New("base_03"),
		},
		Status: true,
	}

	mockLine := line.FromString("mock_line_01")
	unit := New(mock.ToUnit()).
		SetDrawStep(func(_ winsize.Winsize, _ drawable.Unit) ([]line.Line, bool) {
			return []line.Line{mockLine}, false
		}).
		ToUnit()

	unit.Drawable.Boot()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Size(t, 1, lines)
	assert.False(t, status)
	assert.Equal(
		t, text_test.LineToString(mockLine), text_test.LineToString(lines[0]),
	)
}

func TestPipeline_DataStepsChain(t *testing.T) {
	baseLine := line.FromString("base_01")
	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			baseLine,
		},
		Status: true,
	}

	mockLine1 := line.FromString("mock_line_01")
	mockLine2 := line.FromString("mock_line_02")

	unit := New(mock.ToUnit()).
		PushDataSteps(
			func(_ winsize.Winsize, _ drawable.Unit, l []line.Line, s bool) ([]line.Line, bool) {
				return append(l, mockLine1), s
			},
			func(_ winsize.Winsize, _ drawable.Unit, l []line.Line, s bool) ([]line.Line, bool) {
				return append(l, mockLine2), !s
			},
		).ToUnit()

	unit.Drawable.Boot()

	lines, status := unit.Drawable.Draw(winsize.Winsize{})

	assert.Size(t, 3, lines)
	assert.False(t, status)

	assert.Equal(t, text_test.LineToString(baseLine), text_test.LineToString(lines[0]))
	assert.Equal(t, text_test.LineToString(mockLine1), text_test.LineToString(lines[1]))
	assert.Equal(t, text_test.LineToString(mockLine2), text_test.LineToString(lines[2]))
}
