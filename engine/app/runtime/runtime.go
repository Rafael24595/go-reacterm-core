package runtime

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Instance holds the global application Runtime singleton initialized at startup.
var Instance *Runtime

// Runtime manages global execution context state such as session tracking and startup timestamp.
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

// SessionID returns the unique identifier for the current runtime session.
func (r *Runtime) SessionID() string {
	return r.sessionID
}

// Timestamp returns the Unix timestamp (in milliseconds) recorded when the runtime was initialized.
func (r *Runtime) Timestamp() int64 {
	return r.timestamp
}
