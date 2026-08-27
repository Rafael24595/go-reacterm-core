package pager

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestStrategy_DefaultsAndSetters(t *testing.T) {
	strategy := NewStrategy()

	assert.Equal(t, rule.KindPage, strategy.Rule.Kind)
	assert.Equal(t, step.KindPage, strategy.Step.Kind)

	customRule := rule.OnFocus()
	customStep := step.ByLine()

	strategy.WithRule(customRule).WithStep(customStep)

	assert.Equal(t, rule.KindFocus, strategy.Rule.Kind)
	assert.Equal(t, step.KindLine, strategy.Step.Kind)
}

func TestPagerStrategy_Integration(t *testing.T) {
	strategy := NewStrategy()

	ctx := rule.Context{
		Page:     1,
		HasFocus: false,
	}

	pager := state.PagerContext{
		TargetPage: 1,
	}

	shouldStop := strategy.Rule.Handler(pager, ctx)

	assert.True(t, shouldStop)
}
