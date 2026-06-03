package tick

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/behavior"
)

const Tag = "behavior:tick"

func Apply(node screen.Node, decorator behavior.Tick) screen.Node {
	return behavior.Apply(
		node, Wrap(decorator),
	)
}

func Wrap(decorator behavior.Tick) behavior.Behavior {
	return func(node screen.Node) screen.Node {
		node.Screen.Tick = decorator(
			behavior.TargetOf(node),
			node.Screen.Tick,
		)

		node.Tags.Add(Tag)
		
		return node
	}
}
