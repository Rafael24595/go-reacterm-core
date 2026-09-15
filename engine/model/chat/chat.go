package chat

import (
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Message struct {
	Time    int64
	Owner   string
	Message string
}

func NewMessage(owner, content string) Message {
	return NewMessageWithClock(
		owner, content, clock.UnixMilliClock,
	)
}

func NewMessageWithClock(
	owner, content string,
	clock clock.Clock,
) Message {
	return NewMessageWithTimestamp(
		owner, content, clock(),
	)
}

func NewMessageWithTimestamp(
	owner, content string,
	timestamp int64,
) Message {
	return Message{
		Time:    timestamp,
		Owner:   owner,
		Message: content,
	}
}
