package composer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	
	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func Test_PagerRenderer_StaticLayerDoesNotScroll(t *testing.T) {
	uiState := state.NewUIState()
	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetEngine(pager.EngineScroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	dynamic := drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("dyn-1"),
			*text.NewLine("dyn-2"),
			*text.NewLine("dyn-3"),
			*text.NewLine("dyn-4"),
			*text.NewLine("dyn-5"),
			*text.NewLine("dyn-6"),
		},
		Batch: 2,
	}

	static := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("stc"),
		},
		Status: true,
	}

	unit := stack.NewVStack().
		PushLayer(
			dynamic.ToUnit(),
			layer.Fixed[winsize.Rows](2),
		).
		Push(
			static.ToUnit(),
		).
		SetRenderer(renderer).
		ToUnit()

	unit.Drawable.Init()

	size := winsize.New(3, 20)

	page1, hasNext1 := unit.Drawable.Draw(size)
	last1 := page1[len(page1)-1]

	assert.True(t, hasNext1)
	assert.Equal(t, "stc", text.LineToString(&last1))

	page2, _ := unit.Drawable.Draw(size)
	last2 := page2[len(page2)-1]

	assert.Equal(t, text.LineToString(&last1), text.LineToString(&last2))
}

func Test_PagerRenderer_PropagatesMaxPage(t *testing.T) {
	uiState := state.NewUIState()
	uiState.Pager.TargetPage = 3

	ctx := newRenderContext()
	strategy := pager.NewStrategy().
		SetEngine(pager.EngineScroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	mock := drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("1"),
			*text.NewLine("2"),
			*text.NewLine("3"),
			*text.NewLine("4"),
			*text.NewLine("5"),
			*text.NewLine("6"),
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
		SetEngine(pager.EngineScroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("1"),
			*text.NewLine("2"),
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
		SetEngine(pager.EngineScroll())

	renderer := pagerRenderer(uiState, *strategy, ctx)

	unit := (&drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("1"),
			*text.NewLine("2"),
			*text.NewLine("3"),
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
