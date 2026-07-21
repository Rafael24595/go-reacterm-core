package inputline

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestInputLine_UnitBasicSuite(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	unit := New(mock.ToUnit()).ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestNewInputLine_DefaultPrompt(t *testing.T) {
	mock := &drawable_test.MockUnit{}
	input := New(mock.ToUnit())

	assert.Equal(t, input.prompt, marker.DefaultPromptText)
}

func TestNewInputLine_NoContent_ReturnsPromptOnly(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines:  make([]line.Line, 0),
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Boot()
	lines, status := unit.Drawable.Draw(winsize.Winsize{
		Rows: 5,
	})

	assert.False(t, status)
	assert.Size(t, 1, lines)
	assert.Equal(t, marker.DefaultPromptText, text_test.LineToString(lines[0]))
}

func TestNewInputLine_WithSingleLine_AddsPrompt(t *testing.T) {
	frg := frag.FromStrings("golang")

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromFrags(frg...),
		},
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Boot()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Cols: 10,
		Rows: 5,
	})

	assert.Size(t, 1, lines)
	assert.Equal(t, marker.DefaultPromptText+" golang", text_test.LineToString(lines[0]))
}

func TestNewInputLine_MultipleDrawCalls_AccumulatesLines(t *testing.T) {
	frg1 := frag.FromStrings("ziglang")
	frg2 := frag.FromStrings("golang")

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromFrags(frg1...),
			line.FromFrags(frg2...),
		},
	}

	unit := New(mock.ToUnit()).ToUnit()

	unit.Drawable.Boot()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Cols: 10,
		Rows: 5,
	})

	assert.Size(t, 2, lines)

	assert.Equal(t, marker.DefaultPromptText+" ziglang", text_test.LineToString(lines[0]))
	assert.Equal(t, "golang", text_test.LineToString(lines[1]))
}
