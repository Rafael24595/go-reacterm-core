package screen

import "github.com/Rafael24595/go-reacterm-core/engine/model/key"

// Event represents an input event dispatched to the screen's lifecycle,
// encapsulating keyboard inputs such as key actions, modifier masks, and runes.
type Event struct {
	// Key represents the specific keyboard input associated with this event.
	Key key.Key
}

// NewEvent creates and returns a new Event initialized with the provided Key.
func NewEvent(key key.Key) Event {
	return Event{
		Key: key,
	}
}
