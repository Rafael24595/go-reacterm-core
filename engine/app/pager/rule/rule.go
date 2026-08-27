package rule

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

// Kind specifies the type of evaluation rule.
type Kind uint16

const (
	// KindPage evaluates whether the current page matches the target page.
	KindPage Kind = iota
	// KindFocus evaluates whether the current context has focus.
	KindFocus
)

// Context contains the state required to evaluate a rule.
type Context struct {
	Page     uint
	HasFocus bool
}

// Handler is a function that evaluates whether a condition is met.
type Handler func(state.PagerContext, Context) bool

// Rule pairs a classification Kind with an evaluation Handler.
type Rule struct {
	Kind    Kind
	Handler Handler
}

// OnPage evaluates whether the current context page matches the target page.
func OnPage() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(pager state.PagerContext, ctx Context) bool {
			return ctx.Page == pager.TargetPage
		},
	}
}

// OnStart evaluates to always render from the top/beginning boundary.
func OnStart() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return true
		},
	}
}

// OnEnd evaluates relative to the bottom/end boundary.
func OnEnd() Rule {
	return Rule{
		Kind: KindPage,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return false
		},
	}
}

// OnFocus evaluates whether the current context holds active UI focus.
func OnFocus() Rule {
	return Rule{
		Kind: KindFocus,
		Handler: func(_ state.PagerContext, ctx Context) bool {
			return ctx.HasFocus
		},
	}
}
