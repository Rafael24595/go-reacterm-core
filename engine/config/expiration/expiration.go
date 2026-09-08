package expiration

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Expiration struct {
	strategy func(node *screen.Node) bool
}

func Persistent() Expiration {
	return Expiration{}
}

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

func (e Expiration) On(node *screen.Node) bool {
	if e.strategy == nil {
		return false
	}
	return e.strategy(node)
}
