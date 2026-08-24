package screen

import (
	"iter"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

// TODO: Keep it functional?

// Pass defines a compiler pass function transformation applied to a Node.
type Pass func(Node) (Node, error)

// TODO: Hide Stack and build it by adding children data?

// Node represents a component in the UI tree, holding screen logic, lifecycle metadata,
// and hierarchical relations.
type Node struct {
	// id unique identifier generated for the node based on its name and creation sequence.
	id string
	// Name user-defined identifier for the component.
	Name string
	// Tags user or system assigned tags used for categorization and metadata transmission.
	Tags set.Set[string]
	// Screen holds the lifecycle functions (Boot, Keys, Tick, View) for this node.
	Screen Screen
	// Stack specifies the elements under this node to manage flow memory usage.
	Stack set.Set[string]
	// children contains the child nodes attached to this node.
	children []Node
}

// Id returns the unique identifier of the node.
func (n *Node) Id() string {
	return n.id
}

// Children returns an iterator sequence (iter.Seq) over pointers to the node's children.
func (n *Node) Children() iter.Seq[*Node] {
	return func(yield func(*Node) bool) {
		for i := range n.children {
			if !yield(&n.children[i]) {
				return
			}
		}
	}
}

// CompileNode runs a pipeline of Pass transformations sequentially on the node.
// It returns the transformed Node or stops early if a Pass returns an error.
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

// IsZeroNode returns true if any required property of the Node is uninitialized or invalid,
// including whether its Screen is a zero screen.
func IsZeroNode(node Node) bool {
	return node.Name == "" ||
		node.Tags == nil ||
		node.Stack == nil ||
		IsZeroScreen(node.Screen)
}
