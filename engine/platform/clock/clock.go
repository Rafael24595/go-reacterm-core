package clock

import (
	"sync/atomic"
	"time"
)

var globalCounter atomic.Int64

// Clock defines a function provider returning a monotonic integer or timestamp value.
type Clock func() int64

// UnixMilliClock returns the current Unix timestamp in milliseconds.
func UnixMilliClock() int64 {
	return time.Now().UnixMilli()
}

// GlobalCounterClock increments and returns an atomic monotonic counter.
func GlobalCounterClock() int64 {
	return globalCounter.Add(1)

}

func resetGlobalCounter() {
	globalCounter.Store(0)
}
