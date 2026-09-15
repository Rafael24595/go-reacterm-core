package chat_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/chat"
)

func TestNewMessageWithTimestamp(t *testing.T) {
	time := int64(1700000000000)
	owner := "system"
	content := "Hello Golang"

	msg := chat.NewMessageWithTimestamp(owner, content, time)

	assert.Equal(t, time, msg.Timestamp())
	assert.Equal(t, owner, msg.Owner())
	assert.Equal(t, content, msg.Content())
}

func TestNewMessageWithClock(t *testing.T) {
	time := int64(123456789)
	mockClock := func() int64 {
		return time
	}

	msg := chat.NewMessageWithClock("bot", "Automated msg", mockClock)

	assert.Equal(t, time, msg.Timestamp())
}

func TestValueReceiverChaining(t *testing.T) {
	content := "inline call"

	msg := chat.NewMessageWithTimestamp("user", content, 100)

	assert.Equal(t, content, msg.Content())
}
