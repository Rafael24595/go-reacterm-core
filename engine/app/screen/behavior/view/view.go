package view

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
)

const Tag = "behavior:keys"

type Middleware func(uiState state.UIState, context behavior.Context[screen.ViewFunc]) viewmodel.ViewModel

func Apply(node screen.Node, decorator behavior.View) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.View) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.View = decorator(
			behavior.TargetOf(node),
			node.Screen.View,
		)

		node.Tags.Add(Tag)
		return node
	}
}

func Use(node screen.Node, middleware Middleware) screen.Node {
	return Apply(node, use(middleware))
}

func use(middleware Middleware) behavior.View {
	return func(target behavior.Target, next screen.ViewFunc) screen.ViewFunc {
		context := behavior.NewContext(target, next)
		return func(uiState state.UIState) viewmodel.ViewModel {
			return middleware(uiState, context)
		}
	}
}
