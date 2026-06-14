package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

type indexMeta struct {
	sufix      string
	prefixBody string
	digits     uint16
	totalWidth winsize.Cols
}

func (i indexMeta) header(index int) string {
	right := helper.Right(
		winsize.Cols(i.digits),
		helper.TextFromAny(index),
		marker.DefaultPaddingText,
	)

	return right + i.sufix
}

func (i indexMeta) body() string {
	return i.prefixBody + i.sufix
}
