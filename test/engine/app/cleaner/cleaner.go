package cleaner_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/cleaner"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func DiscardCleaner() cleaner.Cleaner {
	return func(res screen.Result, s *state.UIState) *state.UIState {
		return s
	}
}

