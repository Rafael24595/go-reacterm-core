package init

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

const Tag = "behavior:init"

type Middleware func(uiState state.UIState, context behavior.Context[screen.InitFunc])

func Apply(node screen.Node, decorator behavior.Init) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.Init) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Init = decorator(
			behavior.TargetOf(node),
			node.Screen.Init,
		)

		node.Tags.Add(Tag)
		return node
	}
}

func Use(node screen.Node, middleware Middleware) screen.Node {
	return Apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Init {
	return func(target behavior.Target, next screen.InitFunc) screen.InitFunc {
		context := behavior.NewContext(target, next)
		return func(uiState state.UIState) {
			middleware(uiState, context)
		}
	}
}
