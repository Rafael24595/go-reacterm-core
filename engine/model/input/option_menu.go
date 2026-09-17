package input

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type MenuOptionAction = func() screen.Node

type MenuOption struct {
	Id     string
	Label  frag.Frag
	Action MenuOptionAction
}

func NewMenuOption(id string, option frag.Frag) MenuOption {
	return MenuOption{
		Id:     id,
		Label:  option,
	}
}
func (o MenuOption) WithAction(action MenuOptionAction) MenuOption {
	if action == nil {
		return o
	}

	o.Action = action
	return o
}

func (o MenuOption) Exec() (screen.Node, bool) {
	if o.Action != nil {
		return o.Action(), true
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
		if cacheIndex, ok := cache[o.Id]; ok {
			assert.Unreachable("option id '%s' is duplicated", o.Id)

			o.Id = fmt.Sprintf("%s_%d", o.Id, cacheIndex)
			index = cacheIndex + 1
		}

		cache[options[i].Id] = index
		normalized[i] = o
	}

	return normalized
}
