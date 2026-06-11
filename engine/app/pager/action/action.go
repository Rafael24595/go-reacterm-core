package action

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type Kind uint8

const (
	KindPaged Kind = iota
	KindScroll
)

type Handler func(*draw.State) *draw.State

type Action struct {
	Kind    Kind
	Handler Handler
}

func Paged() Action {
	return Action{
		Kind: KindPaged,
		Handler: func(stt *draw.State) *draw.State {
			stt.Reset()
			stt.Page += 1
			return stt
		},
	}
}

func Scroll() Action {
	return Action{
		Kind: KindScroll,
		Handler: func(stt *draw.State) *draw.State {
			if len(stt.Buffer) == 0 {
				return stt
			}

			copy(stt.Buffer, stt.Buffer[1:])
			stt.Buffer[len(stt.Buffer)-1] = text.Line{}
			stt.Cursor = math.SubClampZero(stt.Cursor, 1)

			stt.Focus = false
			stt.Page += 1

			return stt
		},
	}
}
