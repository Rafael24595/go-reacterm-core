package focus

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/transform/page"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// TODO: Add flag to manage non focus unit?
func DrawTransformer(step step.Step) pipeline.DrawTransformer {
	rule := rule.OnFocus()
	strategy := *pager.NewStrategy().
		WithRule(rule).
		WithStep(step)

	return func(size winsize.Winsize, unit drawable.Unit) ([]line.Line, bool) {
		uiState := state.NewUIState()
		renderer := page.NewRenderer(strategy)
		status := renderer(uiState, size, unit)
		return status.Buffer, false
	}
}
