package checkmenu

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
	"github.com/Rafael24595/go-reacterm-core/engine/model/param"
)

const ArgActiveChecks param.Typed[set.Set[string]] = "check_menu_active"
