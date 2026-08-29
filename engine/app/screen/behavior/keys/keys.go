package keys

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
)

// Tag identifies screen nodes that have been decorated with a key definition behavior.
const Tag = "behavior:keys"

// Handler represents a mapping function that transforms or extends a screen's key Definition.
type Handler func(definition screen.Definition) screen.Definition

// Middleware defines an interceptor function with access to the behavior Context during execution of Keys.
type Middleware func(context behavior.Context[screen.KeysFunc]) screen.Definition

// apply decorates the screen Node's Keys function with the given behavior.Keys decorator.
func apply(node screen.Node, decorator behavior.Keys) screen.Node {
	return behavior.Apply(
		node, wrap(decorator),
	)
}

// wrap creates a behavior.Behavior function that decorates a node's Keys function and tags the node.
func wrap(decorator behavior.Keys) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Keys = decorator(
			behavior.TargetOf(node),
			node.Screen.Keys,
		)

		node.Tags.Add(Tag)
		return node
	}
}

// Map attaches a Definition transformation Handler to the screen Node's Keys method.
func Map(node screen.Node, handler Handler) screen.Node {
	return apply(node, mapp(handler))
}

func mapp(handler Handler) behavior.Keys {
	return func(_ behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		return func() screen.Definition {
			return handler(next())
		}
	}
}

// Use intercepts the screen Node's Keys function using a Middleware wrapper.
func Use(node screen.Node, middleware Middleware) screen.Node {
	return apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Keys {
	return func(target behavior.Target, next screen.KeysFunc) screen.KeysFunc {
		context := behavior.NewContext(target, next)
		return func() screen.Definition {
			return middleware(context)
		}
	}
}
