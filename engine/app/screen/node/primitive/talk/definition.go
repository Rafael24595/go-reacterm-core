package talk

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

var predicates = map[bool]predicate.Predicate{
	false: predicate.Page(),
	true:  predicate.Focus(),
}

var actions = map[bool]action.Action{
	false: action.Scroll(),
	true:  action.Paged(),
}

var read_definition = screen.DefinitionFromActions(
	[]key.Action{
		key.ActionEnter,
	}...,
)

var definitions = map[bool]screen.Definition{
	false: read_definition,
	true:  navigation_definition,
}

var navigation_definition = screen.DefinitionFromActions(
	[]key.Action{
		key.ActionEsc,
		key.ActionArrowLeft,
		key.ActionArrowRight,
		key.ActionArrowUp,
		key.ActionArrowDown,
		key.CustomActionPointer,
	}...,
)
