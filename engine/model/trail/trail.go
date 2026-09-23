package trail

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/stack"
)

// DefaultLimit specifies the default history capacity when limit is set to zero.
const DefaultLimit = 3

// Trail maintains back and forward navigation history for screen nodes.
type Trail struct {
	prev    *stack.Stack[screen.Node]
	next    *stack.Stack[screen.Node]
	current screen.Node
}

// New initializes a Trail instance with the specified history limit and starting screen node.
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

// Current returns the active screen node.
func (t *Trail) Current() screen.Node {
	return t.current
}

// Snapshot captures the current state of previous, current, and next navigation history.
func (t *Trail) Snapshot() Snapshot {
	return Snapshot{
		Previous: t.prev.Items(),
		Current:  t.current,
		Next:     t.next.Items(),
	}
}

// GoTo transitions to a new screen node, pushing the current node to previous history
// and clearing any forward navigation history.
func (t *Trail) GoTo(node screen.Node) {
	t.prev.Push(t.current)
	t.current = node
	t.next.Clear()
}

// PeekBack returns the most recent node in the back history without navigating to it.
func (t *Trail) PeekBack() (screen.Node, bool) {
	return t.prev.Peek()
}

// Back navigates to the previous screen node in history if available.
func (t *Trail) Back() (screen.Node, bool) {
	node, ok := t.prev.Pop()
	if !ok {
		var zero screen.Node
		return zero, false
	}

	t.next.Push(t.current)
	t.current = node

	return node, true
}

// PeekForward returns the next node in forward history without navigating to it.
func (t *Trail) PeekForward() (screen.Node, bool) {
	return t.next.Peek()
}

// Forward navigates to the next screen node in forward history if available.
func (t *Trail) Forward() (screen.Node, bool) {
	node, ok := t.next.Pop()
	if !ok {
		var zero screen.Node
		return zero, false
	}

	t.prev.Push(t.current)
	t.current = node

	return node, true
}
