package step

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/draw"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// Kind specifies the type of stepping transition.
type Kind uint8

const (
	// KindPage represents a page-based stepping transition.
	KindPage Kind = iota
	// KindLine represents a scroll-based stepping transition.
	KindLine
)

// Handler modifies the drawing state during a screen step.
type Handler func(*draw.State) *draw.State

// Step holds a Kind identifier and its execution Handler.
type Step struct {
	Kind    Kind
	Handler Handler
}

// ByPage performs a full-page step by resetting state and incrementing the page.
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

// ByLine performs a line shift step by discarding the top buffer line.
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
