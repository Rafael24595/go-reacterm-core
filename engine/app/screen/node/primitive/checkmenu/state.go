package checkmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const KeyActive store.Key[set.Set[string]] = "check_menu_active"
