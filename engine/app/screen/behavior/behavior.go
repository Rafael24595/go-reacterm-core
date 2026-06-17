package behavior

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
)

type Behavior func(screen.Node) screen.Node

type Boot func(target Target, next screen.BootFunc) screen.BootFunc
type Tick func(target Target, next screen.TickFunc) screen.TickFunc
type Keys func(target Target, next screen.KeysFunc) screen.KeysFunc
type View func(target Target, next screen.ViewFunc) screen.ViewFunc

func Apply(node screen.Node, behaviors ...Behavior) screen.Node {
	for _, b := range behaviors {
		node = b(node)
	}
	return node
}
