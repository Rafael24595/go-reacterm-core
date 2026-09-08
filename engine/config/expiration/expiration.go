package expiration

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

// Expiration defines conditions under which a resource or state should expire.
type Expiration struct {
	strategy func(node *screen.Node) bool
}

// Custom creates an Expiration strategy driven by a custom evaluation function.
func Custom(fn func(node *screen.Node) bool) Expiration {
	return Expiration{strategy: fn}
}

// Persistent creates an Expiration strategy that never expires regardless of screen transitions.
func Persistent() Expiration {
	return Expiration{}
}

// OnNode creates an Expiration strategy bound to a specific screen node instance.
// It expires whenever navigation moves away from the specified node.
func OnNode(node *screen.Node) Expiration {
	return Expiration{
		func(next *screen.Node) bool {
			if node != nil {
				return node != next
			}
			return false
		},
	}
}

// OnName creates an Expiration strategy bound to a screen name.
// It triggers expiration when active on a node matching the specified name.
func OnName(name string) Expiration {
	return Expiration{
		strategy: func(next *screen.Node) bool {
			if next != nil {
				return name == next.Name
			}
			return false
		},
	}
}

// On evaluates whether the expiration rule is triggered for the given target screen node.
func (e Expiration) On(node *screen.Node) bool {
	if e.strategy == nil {
		return false
	}
	return e.strategy(node)
}
