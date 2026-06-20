package tick

import (
	"slices"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

const Tag = "behavior:tick"

type Handler func(result screen.Result) screen.Result
type Middleware func(uiState *state.UIState, event screen.Event, context behavior.Context[screen.TickFunc]) screen.Result

func Apply(node screen.Node, decorator behavior.Tick) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.Tick) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Tick = decorator(
			behavior.TargetOf(node),
			node.Screen.Tick,
		)

		node.Tags.Add(Tag)

		return node
	}
}

func Map(node screen.Node, handler Handler) screen.Node {
	return Apply(node, mapp(handler))
}

func mapp(handler Handler) behavior.Tick {
	return func(_ behavior.Target, next screen.TickFunc) screen.TickFunc {
		return func(uiState *state.UIState, event screen.Event) screen.Result {
			return handler(next(uiState, event))
		}
	}
}

func Use(node screen.Node, middleware Middleware) screen.Node {
	return Apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Tick {
	return func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		context := behavior.NewContext(target, next)
		return func(uiState *state.UIState, event screen.Event) screen.Result {
			return middleware(uiState, event, context)
		}
	}
}

func OnKey(node screen.Node, middleware Middleware, keys ...key.Action) screen.Node {
	return Apply(node, onKey(keys, middleware))
}

func onKey(keys []key.Action, middleware Middleware) behavior.Tick {
	return func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		context := behavior.NewContext(target, next)

		return func(uiState *state.UIState, event screen.Event) screen.Result {
			if slices.Contains(keys, event.Key.Code) {
				return middleware(uiState, event, context)
			}
			return next(uiState, event)
		}
	}
}
