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

func CompileNode(n Node, pass ...Pass) (Node, error) {
	node := n

	for _, p := range pass {
		nextNode, err := p(node)
		if err != nil {
			return node, err
		}

		node = nextNode
	}

	return node, nil
}

func IsZeroNode(node Node) bool {
	return node.Name == "" ||
		node.Tags == nil ||
		node.Stack == nil ||
		IsZeroScreen(node.Screen)
}
