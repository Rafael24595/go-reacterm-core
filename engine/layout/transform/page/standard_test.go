package page

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	pager_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/pager"
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestNewPageRenderer_NoEngineCall(t *testing.T) {
	uiState := state.NewUIState()
	size := winsize.Winsize{
		Cols: 10,
		Rows: 5,
	}

	mockStrategy := pager_test.MockStrategy{
		PredicateBool: false,
	}

	strategy := mockStrategy.ToStrategy()

	mock := &drawable_test.MockUnit{
		Lines: make([]line.Line, 1),
	}

	renderer := NewPageRenderer(strategy)
	status := renderer(uiState, size, mock.ToUnit())

	assert.Equal(t, 0, mockStrategy.ActionCall)
	assert.True(t, status.Work.Finished())
	assert.False(t, status.IsFull())
}

func TestNewPageRenderer_EngineCall(t *testing.T) {
	uiState := state.NewUIState()
	size := winsize.Winsize{
		Cols: 10,
		Rows: 5,
	}

	mockStrategy := &pager_test.MockStrategy{
		PredicateBool: false,
		ActionHandler: func(ds *draw.State) *draw.State {
			ds.Reset()
			ds.Page += 1
			return ds
		},
	}

	strategy := mockStrategy.ToStrategy()

	mock := &drawable_test.MockUnit{
		Lines: make([]line.Line, 7),
		Batch: 5,
	}

	renderer := NewPageRenderer(strategy)
	status := renderer(uiState, size, mock.ToUnit())

	assert.Equal(t, 1, mockStrategy.ActionCall)
	assert.Equal(t, 1, status.Page)
	assert.True(t, status.Work.Finished())
	assert.False(t, status.IsFull())
}

func TestNewPageRenderer_EarlyPredicate(t *testing.T) {
	uiState := state.NewUIState()
	size := winsize.Winsize{
		Cols: 10,
		Rows: 5,
	}

	mockStrategy := &pager_test.MockStrategy{
		PredicateBool: true,
	}

	strategy := mockStrategy.ToStrategy()

	mock := &drawable_test.MockUnit{
		Lines: make([]line.Line, 10),
		Batch: 5,
	}

	renderer := NewPageRenderer(strategy)
	status := renderer(uiState, size, mock.ToUnit())

	assert.Equal(t, 0, mockStrategy.ActionCall)
	assert.Equal(t, 1, mockStrategy.PredicateCall)
	assert.True(t, status.Work.Unfinished())
	assert.True(t, status.IsFull())
}

func TestNewPageRenderer_WithLineOverflow(t *testing.T) {
	uiState := state.NewUIState()
	size := winsize.Winsize{
		Cols: 3,
		Rows: 2,
	}

	mockStrategy := &pager_test.MockStrategy{
		ActionHandler: func(ds *draw.State) *draw.State {
			ds.Reset()
			ds.Page += 1
			return ds
		},
		PredicateFunc: func(_ state.PagerContext, pc predicate.Context) bool {
			return pc.Page == 2
		},
	}

	strategy := mockStrategy.ToStrategy()

	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("golang"),
			*line.New("ziglang"),
			*line.New("rust"),
		},
		Batch: 1,
	}

	renderer := NewPageRenderer(strategy)
	status := renderer(uiState, size, mock.ToUnit())

	assert.Equal(t, 2, mockStrategy.ActionCall)
	assert.Equal(t, 2, status.Page)
	assert.True(t, status.Work.Unfinished())
	assert.True(t, status.IsFull())

	assert.Size(t, 2, status.Buffer)

	expected := text_test.LineToString(status.Buffer[0]) +
		text_test.LineToString(status.Buffer[1])
		
	assert.Equal(t, "grus", expected)
}
