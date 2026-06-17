package boot

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

const Tag = "behavior:boot"

type Middleware func(uiState state.UIState, context behavior.Context[screen.BootFunc])

func Apply(node screen.Node, decorator behavior.Boot) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.Boot) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Boot = decorator(
			behavior.TargetOf(node),
			node.Screen.Boot,
		)

		node.Tags.Add(Tag)
		return node
	}
}

func Use(node screen.Node, middleware Middleware) screen.Node {
	return Apply(node, use(middleware))
}

func use(middleware Middleware) behavior.Boot {
	return func(target behavior.Target, next screen.BootFunc) screen.BootFunc {
		context := behavior.NewContext(target, next)
		return func(uiState state.UIState) {
			middleware(uiState, context)
		}
	}
}
