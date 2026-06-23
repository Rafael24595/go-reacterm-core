package pagination

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
)

type Command uint8

const (
	CmdNone Command = iota

	CmdPageUp
	CmdPageDown

	CmdNextPage
	CmdPrevPage
)

var Commands = []Command{
	CmdPageUp,
	CmdPageDown,
	CmdPrevPage,
	CmdNextPage,
}

var defaultBaseBindings = keymap.NewBindings[Command]().
	Bind(key.ActionPageUp, CmdPageUp).
	Bind(key.ActionPageDown, CmdPageDown)

var defaultPagerBindings = keymap.NewBindings[Command]().
	Bind(key.ActionArrowLeft, CmdPrevPage, key.NewDescriptor("Prev page", "←")).
	Bind(key.ActionArrowRight, CmdNextPage, key.NewDescriptor("Next page", "→"))

var defaultScrollBindings = keymap.NewBindings[Command]().
	Bind(key.ActionArrowUp, CmdPrevPage, key.NewDescriptor("Scroll up", "↑")).
	Bind(key.ActionArrowDown, CmdNextPage, key.NewDescriptor("Scroll down", "↓"))

type bindings struct {
	base   *keymap.Bindings[Command]
	pager  *keymap.Bindings[Command]
	scroll *keymap.Bindings[Command]
}

func (d bindings) get(kind action.Kind) *keymap.Bindings[Command] {
	var page *keymap.Bindings[Command]

	switch kind {
	case action.KindPaged:
		page = d.pager
	case action.KindScroll:
		page = d.scroll
	default:
		assert.Unreachable("unhandled action definition %d", kind)
		page = d.pager
	}

	return d.base.Overlay(page)
}

var defaultBindings = bindings{
	base:   defaultBaseBindings,
	pager:  defaultPagerBindings,
	scroll: defaultScrollBindings,
}

type definition struct {
	base   screen.Definition
	pager  screen.Definition
	scroll screen.Definition
}

func emptyDefinition() definition {
	return definition{
		base:   screen.EmptyDefinition(),
		pager:  screen.EmptyDefinition(),
		scroll: screen.EmptyDefinition(),
	}
}

func definitionFromBindings(bindings bindings) definition {
	return definition{
		base:   keymap.BindingsToDefinition(bindings.base),
		pager:  keymap.BindingsToDefinition(bindings.pager),
		scroll: keymap.BindingsToDefinition(bindings.scroll),
	}
}

func (d definition) get(kind action.Kind) screen.Definition {
	var page screen.Definition

	switch kind {
	case action.KindPaged:
		page = d.pager
	case action.KindScroll:
		page = d.scroll
	default:
		assert.Unreachable("unhandled action definition %d", kind)
		page = d.pager
	}

	return d.base.Merge(page)
}

var labels = map[action.Kind]string{
	action.KindPaged:  "page",
	action.KindScroll: "scroll",
}
