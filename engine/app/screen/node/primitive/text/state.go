package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
)

const (
	ArgTextInputState param.Typed[State] = "text_input_state"
	ArgTextInputPulse param.Typed[bool]  = "text_input_pulse"
)

type State struct {
	Write  bool
	Buffer []rune
	Caret  *offset.Offset
	Anchor *offset.Offset
}
