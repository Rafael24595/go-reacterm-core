package text

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
)

const (
	KeyState store.Key[State] = "text_input_state"
	KeySync  store.Key[Sync]  = "text_input_sync"
	KeyPulse store.Key[bool]  = "text_input_pulse"
)

type Sync struct {
	Buffer *[]rune
	Caret  *offset.Offset
	Anchor *offset.Offset
}

type State struct {
	WriteMode bool
	Version   uint64
	Buffer    []rune
	Caret     offset.Offset
	Anchor    offset.Offset
}
