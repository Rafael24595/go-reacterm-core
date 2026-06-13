package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
)

const (
	KeyCursor   store.Key[uint16]         = "talk_cursor"
	KeyMessages store.Key[[]chat.Message] = "talk_messages"
)
