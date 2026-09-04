package store

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestScope_TimestampUpdatesOnMutation(t *testing.T) {
	currentTime := int64(1000)
	mockClock := func() int64 {
		return currentTime
	}

	s := newScope(mockClock)
	assert.Equal(t, 1000, s.timestamp)

	currentTime = 2000
	s.Push("key", newEntry(mockClock, "val"))
	assert.Equal(t, 2000, s.timestamp)

	currentTime = 3000
	s.Remove("key")
	assert.Equal(t, 3000, s.timestamp)
}

func TestRemove_GenericType(t *testing.T) {
	store := New()
	var key Key[string] = "user_role"

	Push(store, "AuthScope", key, "admin")

	val, ok := Remove(store, "AuthScope", key)
	assert.True(t, ok)
	assert.Equal(t, "admin", val)

	_, found := Find(store, "AuthScope", key)
	assert.False(t, found)

	varMissing, okMissing := Remove(store, "AuthScope", key)
	assert.False(t, okMissing)
	assert.Equal(t, "", varMissing)
}

func TestStore_TypeMismatchSafety(t *testing.T) {
	store := New()

	store.Push("Scope", "id", 123)

	var stringKey Key[string] = "id"
	val, ok := Find(store, "Scope", stringKey)

	assert.False(t, ok)
	assert.Equal(t, "", val)

	_, updated := Update(store, "Scope", stringKey, func(s *string) {
		*s = "new_val"
	})
	assert.False(t, updated)
}
