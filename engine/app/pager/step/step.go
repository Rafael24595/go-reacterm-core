package step

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Kind uint8

const (
	KindPage Kind = iota
	KindLine
)

type Handler func(*draw.State) *draw.State

type Step struct {
	Kind    Kind
	Handler Handler
}

func ByPage() Step {
	return Step{
		Kind: KindPage,
		Handler: func(stt *draw.State) *draw.State {
			stt.Reset()
			stt.Page += 1
			return stt
		},
	}
}

func ByLine() Step {
	return Step{
		Kind: KindLine,
		Handler: func(stt *draw.State) *draw.State {
			if len(stt.Buffer) == 0 {
				return stt
			}

			copy(stt.Buffer, stt.Buffer[1:])
			stt.Buffer[len(stt.Buffer)-1] = line.Line{}
			stt.Cursor = math.SubClampZero(stt.Cursor, 1)

			stt.Focus = false
			stt.Page += 1

			return stt
		},
	}
}
