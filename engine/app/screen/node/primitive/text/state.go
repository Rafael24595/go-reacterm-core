package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

const (
	KeyState store.Key[State] = "text_input_state"
	KeyPulse store.Key[bool]  = "text_input_pulse"
)

type State struct {
	Write  bool
	Buffer []rune
	Caret  *offset.Offset
	Anchor *offset.Offset
}
