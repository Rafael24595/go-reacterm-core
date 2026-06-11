package predicate

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

type Kind uint16

const (
	KindPage Kind = iota
	KindFocus
)

type Context struct {
	Page     uint
	HasFocus bool
}

type Handler func(state.PagerContext, Context) bool

type Predicate struct {
	Kind    Kind
	Handler Handler
}

func Page() Predicate {
	return Predicate{
		Kind: KindPage,
		Handler: func(pager state.PagerContext, ctx Context) bool {
			return ctx.Page == pager.TargetPage
		},
	}
}

func Focus() Predicate {
	return Predicate{
		Kind: KindFocus,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return ctx.HasFocus
		},
	}
}
