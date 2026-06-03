package history

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var definition = screen.DefinitionFromActions(
	[]key.Action{
		key.CustomActionBack,
	}...,
)
