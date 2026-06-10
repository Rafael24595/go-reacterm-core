package margin

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
)

func HorizontalFactor(position style.HorizontalPosition) winsize.Cols {
	if position == style.Center {
		return 2
	}
	return 1
}

func VerticalFactor(position style.VerticalPosition) winsize.Rows {
	if position == style.Middle {
		return 2
	}
	return 1
}
