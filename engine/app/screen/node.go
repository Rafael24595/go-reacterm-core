package screen

import "github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"

// TODO: Keep it functional?
type Pass func(Node) (Node, error)

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
	node := n

	for _, m := range pass {
		nextNode, err := m(node)
		if err != nil {
			return node, err
		}

		node = nextNode
	}

	return node, nil
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
