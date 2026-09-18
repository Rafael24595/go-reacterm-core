package offset

import "github.com/Rafael24595/go-reacterm-core/engine/helper/math"

type Offset uint32

func (o Offset) Add(of Offset) Offset {
	return math.AddClampMax(o, of)
}

func (o Offset) Sub(of Offset) Offset {
	return math.SubClampZero(o, of)
}
