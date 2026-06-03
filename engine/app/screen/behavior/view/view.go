package view

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
)

const Tag = "behavior:keys"

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
