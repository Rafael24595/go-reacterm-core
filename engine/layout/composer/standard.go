package composer

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func Standard(
	uiState *state.UIState,
	vm viewmodel.ViewModel,
	size winsize.Winsize,
) (*state.UIState, []line.Line) {
	header := vm.Header.ToUnit()
	header.Drawable.Boot()
	headerLines := drain.UnitEager(size, header)

	footer := vm.Footer.ToUnit()
	footer.Drawable.Boot()
	footerLines := drain.UnitEager(size, footer)

	staticRows := winsize.Rows(
		len(headerLines) + len(footerLines),
	)

	if staticRows > size.Rows {
		return uiState, []line.Line{
			*line.New("Too low resolution"),
		}
	}

	ctx := newRenderContext()

	renderer := pagerRenderer(uiState, *vm.Pager, ctx)

	kernel := vm.Kernel.
		SetRenderer(renderer).
		ToUnit()

	kernel.Drawable.Boot()

	dynamicSize := winsize.New(
		size.Rows.Sub(staticRows),
		size.Cols,
	)

	kernelLines := make([]line.Line, dynamicSize.Rows)

	renderedLines, _ := kernel.Drawable.Draw(dynamicSize)
	copy(kernelLines, renderedLines)

	uiState = syncUIState(uiState, ctx)

	lines := headerLines
	lines = append(lines, kernelLines...)
	lines = append(lines, footerLines...)

	return uiState, lines
}

func syncUIState(uiState *state.UIState, ctx *renderContext) *state.UIState {
	uiState.Pager.ConfirmPage(ctx.MaxPage)
	uiState.Pager.HasMore = ctx.HasMore
	return uiState
}
