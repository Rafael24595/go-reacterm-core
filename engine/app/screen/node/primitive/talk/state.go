package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
)

const (
	KeyState store.Key[State] = "talk_state"
	KeySync  store.Key[Sync]  = "talk_sync"
)

type State struct {
	Cursor   uint16
	Messages []chat.Message
}

type Sync struct {
	Cursor   *uint16
	Messages *[]chat.Message
}
