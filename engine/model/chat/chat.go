package chat

import (
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Message struct {
	timestamp int64
	owner     string
	content   string
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
		timestamp: timestamp,
		owner:     owner,
		content:   content,
	}
}

func (m Message) Timestamp() int64 {
	return m.timestamp
}

func (m Message) Owner() string {
	return m.owner
}

func (m Message) Content() string {
	return m.content
}
