package clip

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
)

const (
	KeyActive  store.Key[bool] = "clip_active"
	KeyRestart store.Key[bool] = "clip_restart"
)
