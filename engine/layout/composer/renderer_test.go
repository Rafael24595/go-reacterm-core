package composer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func Test_PagerRenderer_StaticLayerDoesNotScroll(t *testing.T) {
	uiState := state.NewUIState()
	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetAction(action.Scroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	dynamic := drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("dyn-1"),
			*line.New("dyn-2"),
			*line.New("dyn-3"),
			*line.New("dyn-4"),
			*line.New("dyn-5"),
			*line.New("dyn-6"),
		},
		Batch: 2,
	}

	static := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("stc"),
		},
		Status: true,
	}

	unit := stack.NewVStack().
		PushWithOpts(
			dynamic.ToUnit(),
			layer.Fixed[winsize.Rows](2),
		).
		Push(
			static.ToUnit(),
		).
		SetRenderer(renderer).
		ToUnit()

	unit.Drawable.Boot()

	size := winsize.New(3, 20)

	page1, hasNext1 := unit.Drawable.Draw(size)
	last1 := page1[len(page1)-1]

	assert.True(t, hasNext1)
	assert.Equal(t, "stc", text_test.LineToString(last1))

	page2, _ := unit.Drawable.Draw(size)
	last2 := page2[len(page2)-1]

	assert.Equal(
		t, text_test.LineToString(last1), text_test.LineToString(last2),
	)
}

func Test_PagerRenderer_PropagatesMaxPage(t *testing.T) {
	uiState := state.NewUIState()
	uiState.Pager.TargetPage = 3

	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetAction(action.Scroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	mock := drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("1"),
			*line.New("2"),
			*line.New("3"),
			*line.New("4"),
			*line.New("5"),
			*line.New("6"),
		},
		Batch: 2,
	}

	_, _ = renderer(
		winsize.New(1, 20),
		mock.ToUnit(),
	)

	assert.Equal(t, 3, ctx.MaxPage)
}

func Test_PagerRenderer_SetsHasMore(t *testing.T) {
	uiState := state.NewUIState()
	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetAction(action.Scroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("1"),
			*line.New("2"),
		},
		Batch: 1,
	}

	_, _ = renderer(
		winsize.New(1, 20),
		mock.ToUnit(),
	)

	assert.True(t, ctx.HasMore)
}

func Test_Pager_ConfirmPage_UsesMaxPage(t *testing.T) {
	uiState := state.NewUIState()
	uiState.Pager.TargetPage = 3

	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetAction(action.Scroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	unit := (&drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("1"),
			*line.New("2"),
			*line.New("3"),
		},
		Batch: 1,
	}).ToUnit()

	size := winsize.New(1, 20)

	for range 3 {
		_, _ = renderer(size, unit)
	}

	uiState.Pager.ConfirmPage(ctx.MaxPage)

	assert.Equal(t, 2, uiState.Pager.ActualPage)
}
