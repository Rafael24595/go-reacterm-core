package boot

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

// Tag identifies screen nodes that have been decorated with a boot behavior.
const Tag = "behavior:boot"

// Handler represents a side-effect function executed during the screen's Boot step.
type Handler func()

// Middleware defines an interceptor function with access to UIState and the behavior Context during execution of Boot.
type Middleware func(uiState state.UIState, context behavior.Context[screen.BootFunc])

// apply decorates the screen Node's Boot lifecycle function with the given behavior.Boot decorator.
func apply(node screen.Node, decorator behavior.Boot) screen.Node {
	return behavior.Apply(
		node, wrap(decorator),
	)
}

// wrap creates a behavior.Behavior function that decorates a node's Boot function and tags the node.
func wrap(decorator behavior.Boot) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Boot = decorator(
			behavior.TargetOf(node),
			node.Screen.Boot,
		)

		node.Tags.Add(Tag)
		return node
	}
}

// Map attaches a post-execution Handler to the screen Node's Boot process.
func Map(node screen.Node, handler Handler) screen.Node {
	return apply(node, mapp(handler))
}

func mapp(handler Handler) behavior.Boot {
	return func(_ behavior.Target, next screen.BootFunc) screen.BootFunc {
		return func(uiState state.UIState) {
			next(uiState)
			handler()
		}
	}
}

// Use intercepts the screen Node's Boot function using a Middleware wrapper.
func Use(node screen.Node, middleware Middleware) screen.Node {
	return apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Boot {
	return func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		context := behavior.NewContext(target, next)
		return func(uiState state.UIState) {
			middleware(uiState, context)
		}
	}
}
