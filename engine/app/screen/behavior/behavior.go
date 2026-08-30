package behavior

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
)

// Behavior defines a transformation function that decorates or modifies a screen Node.
type Behavior func(screen.Node) screen.Node

// Boot defines an interceptor wrapper for a screen's BootFunc lifecycle function.
type Boot func(target Target, next screen.BootFunc) screen.BootFunc
// Tick defines an interceptor wrapper for a screen's TickFunc update function.
type Tick func(target Target, next screen.TickFunc) screen.TickFunc
// Keys defines an interceptor wrapper for a screen's KeysFunc key input handler.
type Keys func(target Target, next screen.KeysFunc) screen.KeysFunc
// View defines an interceptor wrapper for a screen's ViewFunc rendering function.
type View func(target Target, next screen.ViewFunc) screen.ViewFunc

// Apply sequentially applies a series of Behavior transformations to a screen Node.
func Apply(node screen.Node, behaviors ...Behavior) screen.Node {
	for _, b := range behaviors {
		node = b(node)
	}
	return node
}
