package layout_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type MockLayout struct {
	Composer layout.Composer
}

func FromComposer(composer layout.Composer) layout.Layout {
	return MockLayout{
		Composer: composer,
	}.ToLayout()
}

func DiscardLayout() layout.Layout {
	return MockLayout{}.ToLayout()
}

func (m MockLayout) ToLayout() layout.Layout {
	if m.Composer == nil {
		m.Composer = func(s *state.UIState, vm viewmodel.ViewModel, size winsize.Winsize) (*state.UIState, []line.Line) {
			return s, make([]line.Line, 0)
		}
	}

	return layout.Layout{
		Compose: m.Composer,
	}
}
