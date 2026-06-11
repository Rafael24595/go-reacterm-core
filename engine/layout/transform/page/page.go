package page

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type Renderer func(*state.UIState, winsize.Winsize, drawable.Unit) *draw.State
