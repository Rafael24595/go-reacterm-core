package runtime

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var Instance *Runtime

type Runtime struct {
	sessionID string
	timestamp int64
}

func init() {
	Instance = &Runtime{
		sessionID: newSessionId(),
		timestamp: time.Now().UnixMilli(),
	}
}

func newSessionId() string {
	now := time.Now().UnixNano()
	return fmt.Sprintf("%d-%04x", now, rand.Uint32())
}

func (r *Runtime) SessionID() string {
	return r.sessionID
}

func (r *Runtime) Timestamp() int64 {
	return r.timestamp
}
