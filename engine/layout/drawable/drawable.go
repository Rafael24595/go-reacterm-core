package drawable

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type BootFunc func()
type WipeFunc func()
type DrawFunc func(size winsize.Winsize) ([]line.Line, bool)

type Drawable struct {
	Boot BootFunc
	Wipe WipeFunc
	Draw DrawFunc
}

func IsZeroDrawable(drawable Drawable) bool {
	if drawable.Boot == nil {
		return true
	}

	if drawable.Wipe == nil {
		return true
	}

	if drawable.Draw == nil {
		return true
	}

	return false
}
