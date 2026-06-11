package action

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

func TestActionPaged(t *testing.T) {
	action := Paged()

	state := &draw.State{
		Buffer: []text.Line{{}, {}, {}},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := action.Handler(state)

	assert.Size(t, 3, result.Buffer)
	assert.Equal(t, 2, result.Page)
	assert.Equal(t, 0, result.Cursor)
	assert.False(t, result.Focus)
}

func TestActionPaged_AlwaysResetsBuffer(t *testing.T) {
	action := Paged()

	state := &draw.State{
		Buffer: []text.Line{{}, {}},
	}

	action.Handler(state)
	action.Handler(state)

	assert.Equal(t, 2, state.Page)
}

func TestActionScroll(t *testing.T) {
	action := Scroll()

	state := &draw.State{
		Buffer: []text.Line{
			*text.NewLine("A"),
			*text.NewLine("B"),
			*text.NewLine("C"),
		},
		Cursor: 2,
		Page:   1,
		Focus:  true,
	}

	result := action.Handler(state)

	assert.Equal(t, "B", text.LineToString(&result.Buffer[0]))
	assert.Equal(t, "C", text.LineToString(&result.Buffer[1]))
	assert.Equal(t, "", text.LineToString(&result.Buffer[2]))
	assert.Equal(t, 1, result.Cursor)
	assert.False(t, result.Focus)
}

func TestActionScroll_CursorNeverNegative(t *testing.T) {
	action := Scroll()

	state := &draw.State{
		Cursor: 0,
	}

	result := action.Handler(state)

	assert.Equal(t, 0, result.Cursor)
}
