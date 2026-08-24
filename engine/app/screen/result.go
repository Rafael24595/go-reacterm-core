package screen

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

// Result represents the outcome of a screen's Tick execution, specifying state transitions
// or updates to the navigation/pager context.
type Result struct {
	// node target Node to navigate or update to.
	node    Node
	// hasNode indicates whether a target Node is present in this result.
	hasNode bool
	// Pager holds the pager state context associated with screen navigation operations.
	Pager   state.PagerContext
}

// EmptyResult initializes and returns a Result with default uninitialized values and no target Node.
func EmptyResult() Result {
	return Result{
		node:    Node{},
		hasNode: false,
		Pager:   state.PagerContext{},
	}
}

// ResultFromNode constructs a Result initialized with a target navigation Node.
func ResultFromNode(node Node) Result {
	return Result{
		node:    node,
		hasNode: true,
		Pager:   state.PagerContext{},
	}
}

// ResultFromUIState constructs a Result populated with the PagerContext extracted from a UIState pointer.
func ResultFromUIState(uiState *state.UIState) Result {
	if uiState == nil {
		return EmptyResult()
	}
	
	return Result{
		node:    Node{},
		hasNode: false,
		Pager:   uiState.Pager,
	}
}

// HasNode returns true if the Result contains a valid target Node.
func (r *Result) HasNode() bool {
	return r.hasNode
}

// GetNode returns the target Node held by the Result.
func (r *Result) GetNode() Node {
	return r.node
}

// TryGetNode returns the target Node alongside a boolean indicating whether the Node is present.
func (r *Result) TryGetNode() (Node, bool) {
	return r.node, r.hasNode
}

// SetNode assigns a target Node to the Result and marks hasNode as true, returning a pointer for method chaining.
func (r *Result) SetNode(node Node) *Result {
	r.node = node
	r.hasNode = true
	return r
}
