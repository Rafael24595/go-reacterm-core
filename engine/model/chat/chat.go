package chat

import (
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

// Message represents a single chat payload sent by a user.
type Message struct {
	timestamp int64
	owner     string
	content   string
}

// NewMessage initializes a Message with the current timestamp using the default clock function.
func NewMessage(owner, content string) Message {
	return NewMessageWithClock(
		owner, content, clock.UnixMilliClock,
	)
}

// NewMessageWithClock initializes a Message with a custom clock function for testing or other purposes.
func NewMessageWithClock(
	owner, content string,
	clock clock.Clock,
) Message {
	return NewMessageWithTimestamp(
		owner, content, clock(),
	)
}

// NewMessageWithTimestamp initializes a Message with a specific timestamp, allowing for precise control over the message's creation time.
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

// Timestamp returns the timestamp of the message, indicating when it was created.
func (m Message) Timestamp() int64 {
	return m.timestamp
}

// Owner returns the owner of the message, indicating who sent it.
func (m Message) Owner() string {
	return m.owner
}

// Content returns the content of the message, representing the actual text or payload sent by the user.
func (m Message) Content() string {
	return m.content
}
