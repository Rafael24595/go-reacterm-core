package page

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap"
)

func NewPageRenderer(strategy pager.PagerStrategy) Renderer {
	return func(uiState *state.UIState, size winsize.Winsize, unit drawable.Unit) *draw.State {
		status := draw.NewState(size.Rows)
		if size.Rows == 0 {
			return status
		}

		status.Work.Add(1)

		for status.Work.Unfinished() {
			status.Work.Advance()
			status.Work.Reset()

			lines, hasNext := unit.Drawable.Draw(size)
			if hasNext {
				status.Work.Add(1)
			}

			linesLen := uint(len(lines))
			if linesLen == 0 {
				return status
			}

			status.Work.Add(linesLen)

			for _, lne := range lines {
				fixed := wrap.Line(size.Cols, lne)

				fixedLen := uint(len(fixed))
				if fixedLen == 0 {
					continue
				}

				status.Work.Advance()
				status.Work.Add(fixedLen)

				for _, fix := range fixed {
					status.SetAndNext(fix)
					status.Work.Advance()

					status.MarkFocus(
						line.HasAtom(atom.Focus, fix),
					)

					if !status.IsFull() {
						continue
					}

					if shouldStop(uiState.Pager, strategy, status) {
						return status
					}

					if status.Work.Unfinished() {
						status = strategy.Action.Handler(status)
					}
				}
			}
		}

		return status
	}
}

func shouldStop(
	ctx state.PagerContext,
	strategy pager.PagerStrategy,
	status *draw.State,
) bool {
	args := predicate.Context{
		Page:     status.Page,
		HasFocus: status.Focus,
	}
	return strategy.Predicate.Handler(ctx, args)
}
