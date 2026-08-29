package tick

import (
	"slices"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

// Tag identifies screen nodes that have been decorated with a tick update behavior.
const Tag = "behavior:tick"

// Handler represents a transformation function that processes a screen's execution Result.
type Handler func(result screen.Result) screen.Result

// Middleware defines an interceptor function with access to state, events, and the behavior Context during execution of Tick.
type Middleware func(uiState *state.UIState, event screen.Event, context behavior.Context[screen.TickFunc]) screen.Result

// apply decorates the screen Node's Tick function with the given behavior.Tick decorator.
func apply(node screen.Node, decorator behavior.Tick) screen.Node {
	return behavior.Apply(
		node, wrap(decorator),
	)
}

// wrap creates a behavior.Behavior function that decorates a node's Tick function and tags the node.
func wrap(decorator behavior.Tick) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Tick = decorator(
			behavior.TargetOf(node),
			node.Screen.Tick,
		)

		node.Tags.Add(Tag)

		return node
	}
}

// Map attaches a Result transformation Handler to the screen Node's Tick method.
func Map(node screen.Node, handler Handler) screen.Node {
	return apply(node, mapp(handler))
}

func mapp(handler Handler) behavior.Tick {
	return func(_ behavior.Target, next screen.TickFunc) screen.TickFunc {
		return func(uiState *state.UIState, event screen.Event) screen.Result {
			return handler(next(uiState, event))
		}
	}
}

// Use intercepts the screen Node's Tick function using a Middleware wrapper.
func Use(node screen.Node, middleware Middleware) screen.Node {
	return apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Tick {
	return func(target behavior.Target, next screen.TickFunc) screen.TickFunc {
		context := behavior.NewContext(target, next)
		return func(uiState *state.UIState, event screen.Event) screen.Result {
			return middleware(uiState, event, context)
		}
	}
}

// OnKey restricts the execution of the given Middleware to specific key action triggers matching the input event.
func OnKey(node screen.Node, middleware Middleware, keys ...key.Action) screen.Node {
	return apply(node, onKey(keys, middleware))
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
