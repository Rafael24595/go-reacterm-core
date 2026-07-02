package clip

import (
	"time"

	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
)

const (
	KeyState store.Key[State] = "clip_state"
	KeySync  store.Key[Sync]  = "clip_sync"

	KeyRestart store.Key[bool] = "clip_restart"
)

type Sync struct {
	Active *bool
	Pause  *time.Duration
}

type State struct {
	Active bool
	Pause  time.Duration
}
