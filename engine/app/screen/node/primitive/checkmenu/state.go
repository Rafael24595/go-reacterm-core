package checkmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	KeyState store.Key[set.Set[string]] = "check_menu_state"
	KeySync  store.Key[set.Set[string]] = "check_menu_sync"
)
