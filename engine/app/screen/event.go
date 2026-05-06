package screen

import "github.com/Rafael24595/go-reacterm-core/engine/model/key"

type Event struct {
	Key key.Key
}

func NewEvent(key key.Key) Event {
	return Event{
		Key: key,
	}
}
