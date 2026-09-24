package clock

import (
	"sync/atomic"
	"time"
)

var globalCounter atomic.Int64

type Clock func() int64

func UnixMilliClock() int64 {
	return time.Now().UnixMilli()
}

func GlobalCounterClock() int64 {
	return globalCounter.Add(1)

}

func resetGlobalCounter() {
	globalCounter.Store(0)
}
