package state

import "github.com/Rafael24595/go-reacterm-core/engine/app/store"

// UIState aggregates all sub-contexts representing the visual and operational state of the UI.
type UIState struct {
	Helper HelperContext
	Pager  PagerContext
	Store  *store.Store
}

// NewUIState initializes a UIState with default contexts and a fresh Store instance.
func NewUIState() *UIState {
	return &UIState{
		Helper: HelperContext{},
		Pager:  PagerContext{},
		Store:  store.New(),
	}
}
