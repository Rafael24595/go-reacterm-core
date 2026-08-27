package rule

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

type Rule struct {
	Kind    Kind
	Handler Handler
}

func OnPage() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(pager state.PagerContext, ctx Context) bool {
			return ctx.Page == pager.TargetPage
		},
	}
}

func OnStart() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return true
		},
	}
}

func OnEnd() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return false
		},
	}
}

func OnFocus() Rule {
	return Rule{
		Kind: KindFocus,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return ctx.HasFocus
		},
	}
}
