package pager

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
)

var (
	defaultRule = rule.OnPage()
	defaultStep = step.ByPage()
)

type Strategy struct {
	Rule rule.Rule
	Step step.Step
}

func NewStrategy() *Strategy {
	return &Strategy{
		Rule: defaultRule,
		Step: defaultStep,
	}
}

func (p *Strategy) WithRule(rule rule.Rule) *Strategy {
	p.Rule = rule
	return p
}

func (p *Strategy) WithStep(step step.Step) *Strategy {
	p.Step = step
	return p
}
