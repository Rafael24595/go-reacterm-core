package input

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type MenuOptionHandler = func() screen.Node

type MenuOption struct {
	id      string
	label   frag.Frag
	handler MenuOptionHandler
}

func NewMenuOption(id string, label frag.Frag) MenuOption {
	return MenuOption{
		id:    id,
		label: label,
	}
}

func (o MenuOption) WithHandler(handler MenuOptionHandler) MenuOption {
	if handler == nil {
		return o
	}

	o.handler = handler
	return o
}

func (o MenuOption) Id() string {
	return o.id
}

func (o MenuOption) Label() frag.Frag {
	return o.label
}

func (o MenuOption) Exec() (screen.Node, bool) {
	if o.handler != nil {
		return o.handler(), true
	}
	return screen.Node{}, false
}

func ExtractMenuOptionLabels(options ...MenuOption) []frag.Frag {
	lines := make([]frag.Frag, len(options))
	for i := range options {
		lines[i] = options[i].label
	}
	return lines
}

func NormalizeMenuOptions(options ...MenuOption) []MenuOption {
	normalized := make([]MenuOption, len(options))
	cache := make(map[string]uint)

	for i, o := range options {
		index := uint(1)
		if cacheIndex, ok := cache[o.id]; ok {
			assert.Unreachable("option id '%s' is duplicated", o.id)

			o.id = fmt.Sprintf("%s_%d", o.id, cacheIndex)
			index = cacheIndex + 1
		}

		cache[o.id] = index
		normalized[i] = o
	}

	return normalized
}
