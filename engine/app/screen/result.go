package screen

import "github.com/Rafael24595/go-reacterm-core/engine/app/state"

// TODO: Use node copy instead reference?
type Result struct {
	Node    *Node
	Pager   state.PagerContext
}

func ResultFromNode(node *Node) Result {
	return Result{
		Node:    node,
		Pager:   state.PagerContext{},
	}
}

func ResultFromUIState(uiState *state.UIState) Result {
	return Result{
		Node:    nil,
		Pager:   uiState.Pager,
	}
}

func EmptyResult() Result {
	return Result{
		Node:    nil,
		Pager:   state.PagerContext{},
	}
}
