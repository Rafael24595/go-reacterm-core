package step

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestKinds(t *testing.T) {
	assert.Equal(t, KindPage, ByPage().Kind)
	assert.Equal(t, KindLine, ByLine().Kind)
}

func TestByPage(t *testing.T) {
	step := ByPage()

	state := &draw.State{
		Buffer: []line.Line{{}, {}, {}},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := step.Handler(state)

	assert.Size(t, 3, result.Buffer)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 0, result.Cursor)
	assert.False(t, result.Focus)
}

func TestByPage_AlwaysResetsBuffer(t *testing.T) {
	step := ByPage()

	state := &draw.State{
		Buffer: []line.Line{{}, {}},
	}

	step.Handler(state)
	step.Handler(state)

	assert.Equal(t, 2, state.Page)
}

func TestByLine(t *testing.T) {
	step := ByLine()

	state := &draw.State{
		Buffer: []line.Line{
			line.FromString("A"),
			line.FromString("B"),
			line.FromString("C"),
		},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := step.Handler(state)

	assert.Equal(t, "B", text_test.LineToString(result.Buffer[0]))
	assert.Equal(t, "C", text_test.LineToString(result.Buffer[1]))
	assert.Equal(t, "", text_test.LineToString(result.Buffer[2]))
	assert.Equal(t, 1, result.Cursor)
	assert.False(t, result.Focus)
}

func TestByLine_CursorNeverNegative(t *testing.T) {
	step := ByLine()

	state := &draw.State{
		Cursor: 0,
	}

	result := step.Handler(state)

	assert.Equal(t, 0, result.Cursor)
}

func TestByLine_EmptyBuffer(t *testing.T) {
	step := ByLine()
	state := &draw.State{
		Buffer: []line.Line{},
		Cursor: 0,
		Page:   1,
	}

	result := step.Handler(state)

	assert.Size(t, 0, result.Buffer)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 0, result.Cursor)
}
