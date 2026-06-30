package trail

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/stack"
)

const DefaultLimit = 3

type Trail struct {
	prev    *stack.Stack[screen.Node]
	next    *stack.Stack[screen.Node]
	current screen.Node
}

func New(limit uint, current screen.Node) *Trail {
	if limit == 0 {
		assert.Unreachable("limit should be greater than 0")

		limit = DefaultLimit
	}

	return &Trail{
		prev:    stack.New[screen.Node](limit),
		next:    stack.New[screen.Node](limit),
		current: current,
	}
}

func (t *Trail) Current() screen.Node {
	return t.current
}

func (t *Trail) Snapshot() Snapshot {
	return Snapshot{
		Previous: t.prev.Items(),
		Current:  t.current,
		Next:     t.next.Items(),
	}
}

func (n *Trail) GoTo(node screen.Node) {
	n.prev.Push(n.current)
	n.current = node
	n.next.Clear()
}

func (n *Trail) PeekBack() (screen.Node, bool) {
	return n.prev.Peek()
}

func (n *Trail) Back() (screen.Node, bool) {
	node, ok := n.prev.Pop()
	if !ok {
		var zero screen.Node
		return zero, false
	}

	n.next.Push(n.current)
	n.current = node

	return node, true
}

func (n *Trail) PeekForward() (screen.Node, bool) {
	return n.next.Peek()
}

func (n *Trail) Forward() (screen.Node, bool) {
	node, ok := n.next.Pop()
	if !ok {
		var zero screen.Node
		return zero, false
	}

	n.prev.Push(n.current)
	n.current = node

	return node, true
}
