package pager

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/rule"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/step"
)

var (
	defaultRule = rule.OnPage()
	defaultStep = step.ByPage()
)

// Strategy configures how content is filtered and stepped across pages.
type Strategy struct {
	// Rule defines the condition to determine when to stop paging.
	Rule rule.Rule
	// Step defines how to advance to the next page or scroll.
	Step step.Step
}

// NewStrategy creates a Strategy initialized with default rule and step behaviors.
func NewStrategy() *Strategy {
	return &Strategy{
		Rule: defaultRule,
		Step: defaultStep,
	}
}

// WithRule sets the evaluation rule for the strategy.
func (p *Strategy) WithRule(rule rule.Rule) *Strategy {
	p.Rule = rule
	return p
}

// WithStep sets the stepping behavior for the strategy.
func (p *Strategy) WithStep(step step.Step) *Strategy {
	p.Step = step
	return p
}

// Clone creates a deep copy of the Strategy, preserving its rule and step configurations.
func (p *Strategy) Clone() *Strategy {
	return &Strategy{
		Rule: rule.Rule{
			Kind:    p.Rule.Kind,
			Handler: p.Rule.Handler,
		},
		Step: step.Step{
			Kind:    p.Step.Kind,
			Handler: p.Step.Handler,
		},
	}
}
