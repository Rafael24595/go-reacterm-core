package view

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

// Tag identifies screen nodes that have been decorated with a view behavior.
const Tag = "behavior:keys"

// Handler represents a transformation function that processes a screen's ViewModel.
type Handler func(vm viewmodel.ViewModel) viewmodel.ViewModel

// Middleware defines an interceptor function with access to the UI state and the behavior Context during execution of View.
type Middleware func(uiState state.UIState, context behavior.Context[screen.ViewFunc]) viewmodel.ViewModel

// apply decorates the screen Node's View function with the given behavior.View decorator.
func apply(node screen.Node, decorator behavior.View) screen.Node {
	return behavior.Apply(
		node, wrap(decorator),
	)
}

// wrap creates a behavior.Behavior function that decorates a node's View function and tags the node.
func wrap(decorator behavior.View) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.View = decorator(
			behavior.TargetOf(node),
			node.Screen.View,
		)

		node.Tags.Add(Tag)
		return node
	}
}

// Map attaches a ViewModel transformation Handler to the screen Node's View method.
func Map(node screen.Node, handler Handler) screen.Node {
	return apply(node, mapp(handler))
}

func mapp(handler Handler) behavior.View {
	return func(_ behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		return func(uiState state.UIState) viewmodel.ViewModel {
			return handler(next(uiState))
		}
	}
}

// Use intercepts the screen Node's View function using a Middleware wrapper.
func Use(node screen.Node, middleware Middleware) screen.Node {
	return apply(node, use(middleware))
}

func use(middleware Middleware) behavior.View {
	return func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		context := behavior.NewContext(target, next)
		return func(uiState state.UIState) viewmodel.ViewModel {
			return middleware(uiState, context)
		}
	}
}
