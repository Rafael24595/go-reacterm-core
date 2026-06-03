package keys

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
)

const Tag = "behavior:keys"

func Apply(node screen.Node, decorator behavior.Keys) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.Keys) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Keys = decorator(
			behavior.TargetOf(node),
			node.Screen.Keys,
		)

		node.Tags.Add(Tag)
		return node
	}
}
