package state

import "github.com/Rafael24595/go-reacterm-core/engine/app/store"

type UIState struct {
	Helper HelperContext
	Pager  PagerContext
	Store  *store.Store
}

func NewUIState() *UIState {
	return &UIState{
		Helper: HelperContext{},
		Pager:  PagerContext{},
		Store:  store.New(),
	}
}
