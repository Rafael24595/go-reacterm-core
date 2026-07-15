package lines

import (
	"github.com/Rafael24595/go-reacterm-core/engine/format"
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
	right := format.JustifyLeft(
		winsize.Cols(i.digits),
		format.TextFromAny(index),
		marker.DefaultPaddingText,
	)

	return right + i.sufix
}

func (i indexMeta) body() string {
	return i.prefixBody + i.sufix
}
