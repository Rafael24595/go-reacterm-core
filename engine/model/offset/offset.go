package offset

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type Offset uint32

func (r Offset) Sub(o Offset) Offset {
	return math.SubClampZero(r, o)
}
