package input

import (
	"fmt"

	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

// MenuOptionHandler handles the execution callback of a menu option.
type MenuOptionHandler = func() screen.Node

// MenuOption represents an actionable menu entry with an identifier, label fragment, and handler.
type MenuOption struct {
	id      string
	label   frag.Frag
	handler MenuOptionHandler
}

// NewMenuOption initializes a new MenuOption with an ID and label fragment.
func NewMenuOption(id string, label frag.Frag) MenuOption {
	return MenuOption{
		id:    id,
		label: label,
	}
}

// WithHandler updates the underlying callback handler.
func (o MenuOption) WithHandler(handler MenuOptionHandler) MenuOption {
	if handler == nil {
		return o
	}

	o.handler = handler
	return o
}

// Id returns the option unique identifier.
func (o MenuOption) Id() string {
	return o.id
}

// Label returns the text fragment label of the option.
func (o MenuOption) Label() frag.Frag {
	return o.label
}

// Exec safely invokes the underlying handler.
func (o MenuOption) Exec() (screen.Node, bool) {
	if o.handler != nil {
		return o.handler(), true
	}
	return screen.Node{}, false
}

// ExtractCheckOptionLabels extracts the underlying frag.Frag labels from a collection of MenuOptions.
func ExtractMenuOptionLabels(options ...MenuOption) []frag.Frag {
	lines := make([]frag.Frag, len(options))
	for i := range options {
		lines[i] = options[i].label
	}
	return lines
}

// NormalizeMenuOptions ensures option IDs are unique.
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
