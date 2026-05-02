package line

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

type indexMeta struct {
	sufix      string
	prefixBody string
	digits     uint16
	totalWidth winsize.Cols
}

func (i indexMeta) header(index int) string {
	return helper.Right(index, winsize.Cols(i.digits)) + i.sufix
}

func (i indexMeta) body() string {
	return i.prefixBody + i.sufix
}
