package store

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestEntry_CreationAndAccessors(t *testing.T) {
	currentTime := int64(1257030000)
	mockClock := func() int64 {
		return currentTime
	}

	entry := newEntry(mockClock, "golang")

	assert.Equal(t, currentTime, entry.Timestamp())
	assert.Equal(t, "golang", entry.Value().Text())
}
