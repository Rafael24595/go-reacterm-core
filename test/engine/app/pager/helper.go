package pager_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

type MockStrategy struct {
	ActionCall    uint
	ActionKind    action.Kind
	ActionHandler action.Handler
	PredicateCall uint
	PredicateKind predicate.Kind
	PredicateBool bool
	PredicateFunc predicate.Handler
}

func (s *MockStrategy) ToStrategy() pager.PagerStrategy {
	return pager.PagerStrategy{
		Action: action.Action{
			Kind: s.ActionKind,
			Handler: func(ds *draw.State) *draw.State {
				s.ActionCall += 1
				if s.ActionHandler != nil {
					return s.ActionHandler(ds)
				}
				return ds
			},
		},
		Predicate: predicate.Predicate{
			Kind: s.PredicateKind,
			Handler: func(c state.PagerContext, pc predicate.Context) bool {
				s.PredicateCall += 1
				if s.PredicateFunc != nil {
					return s.PredicateFunc(c, pc)
				}
				return s.PredicateBool
			},
		},
	}
}
