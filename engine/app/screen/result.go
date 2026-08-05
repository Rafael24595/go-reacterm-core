package screen

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

type Result struct {
	node    Node
	hasNode bool
	Pager   state.PagerContext
}

func ResultFromNode(node Node) Result {
	return Result{
		node:    node,
		hasNode: true,
		Pager:   state.PagerContext{},
	}
}

func ResultFromUIState(uiState *state.UIState) Result {
	return Result{
		node:    Node{},
		hasNode: false,
		Pager:   uiState.Pager,
	}
}

func EmptyResult() Result {
	return Result{
		node:    Node{},
		hasNode: false,
		Pager:   state.PagerContext{},
	}
}

func (r *Result) HasNode() bool {
	return r.hasNode
}

func (r *Result) TryGetNode() (Node, bool) {
	return r.node, r.hasNode
}


func (r *Result) GetNode() Node {
	return r.node
}

func (r *Result) SetNode(node Node) *Result {
	r.node = node
	r.hasNode = true
	return r
}
