package pager_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

type MockStrategy struct {
	StepCall    uint
	StepKind    step.Kind
	StepHandler step.Handler

	RuleCall    uint
	RuleKind    rule.Kind
	RuleBool    bool
	RuleHandler rule.Handler
}

func (s *MockStrategy) ToStrategy() pager.Strategy {
	return pager.Strategy{
		Step: step.Step{
			Kind: s.StepKind,
			Handler: func(ds *draw.State) *draw.State {
				s.StepCall += 1
				if s.StepHandler != nil {
					return s.StepHandler(ds)
				}
				return ds
			},
		},
		Rule: rule.Rule{
			Kind: s.RuleKind,
			Handler: func(c state.PagerContext, ctx rule.Context) bool {
				s.RuleCall += 1
				if s.RuleHandler != nil {
					return s.RuleHandler(c, ctx)
				}
				return s.RuleBool
			},
		},
	}
}
