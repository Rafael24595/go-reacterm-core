package input

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type MenuOptionAction = func() screen.Node

type MenuOption struct {
	Id     string
	Label  text.Fragment
	Action MenuOptionAction
}

func NewMenuOption(id string, option text.Fragment, action MenuOptionAction) MenuOption {
	return MenuOption{
		Id:     id,
		Label:  option,
		Action: action,
	}
}

func NewMenuOptions(options ...MenuOption) []MenuOption {
	return options
}

func FragmentFromMenuOption(options ...MenuOption) []text.Fragment {
	lines := make([]text.Fragment, len(options))
	for i := range options {
		lines[i] = options[i].Label
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
