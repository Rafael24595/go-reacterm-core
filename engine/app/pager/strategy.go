package pager

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
)

var (
	default_predicate = predicate.Page()
	default_action    = action.Paged()
)

type PagerStrategy struct {
	Predicate predicate.Predicate
	Action    action.Action
}

func NewStrategy() *PagerStrategy {
	return &PagerStrategy{
		Predicate: default_predicate,
		Action:    default_action,
	}
}

func (p *PagerStrategy) SetPredicate(predicate predicate.Predicate) *PagerStrategy {
	p.Predicate = predicate
	return p
}

func (p *PagerStrategy) SetAction(action action.Action) *PagerStrategy {
	p.Action = action
	return p
}
