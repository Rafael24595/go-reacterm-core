package screen

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

type Node struct {
	id       string
	Name     string
	Tags     set.Set[string]
	Screen   Screen
	Stack    set.Set[string]
	children []Node
}

func (n Node) Id() string {
	return n.id
}

func (n Node) Children() []Node {
	return n.children
}

func (n Node) Compile(pass ...Pass) (Node, error) {
	screen := n

	for _, m := range pass {
		nextScreen, err := m(screen)
		if err != nil {
			return screen, err
		}

		screen = nextScreen
	}

	return screen, nil
}

func IsZeroNode(node Node) bool {
	if node.Name == "" {
		return true
	}

	if node.Tags == nil {
		return true
	}

	if IsZeroScreen(node.Screen) {
		return true
	}

	if node.Stack == nil {
		return true
	}

	return false
}
