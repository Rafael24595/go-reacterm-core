package store

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/dynamic"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
)

type Entry struct {
	timestamp int64
	value     dynamic.Value
}

func newEntry(clk clock.Clock, arg any) Entry {
	return Entry{
		timestamp: clk(),
		value:     dynamic.From(arg),
	}
}

func (e Entry) Timestamp() int64 {
	return e.timestamp
}

func (e Entry) Value() dynamic.Value {
	return e.value
}
