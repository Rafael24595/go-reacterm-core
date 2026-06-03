package init

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
)

const Tag = "behavior:init"

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
