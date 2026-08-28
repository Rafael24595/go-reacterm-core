package pass

import (
	"errors"
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

const (
	err_name   = "screen: name is required"
	errf_boot  = "screen %q: Boot is nil"
	errf_keys  = "screen %q: Keys is nil"
	errf_tick  = "screen %q: Tick is nil"
	errf_view  = "screen %q: View is nil"
	errf_cycle = "screen %q: Cycle detected"
)

func ValidateStructure(node screen.Node) (screen.Node, error) {
	visited := set.New[string]()

	pending := []screen.Node{node}
	cursor := 0

	for cursor < len(pending) {
		focus := pending[cursor]
		visited.Add(focus.Id())

		if focus.Name == "" {
			return node, errors.New(err_name)
		}

		if focus.Screen.Boot == nil {
			return node, fmt.Errorf(errf_boot, focus.Name)
		}

		if focus.Screen.Keys == nil {
			return node, fmt.Errorf(errf_keys, focus.Name)
		}

		if focus.Screen.Tick == nil {
			return node, fmt.Errorf(errf_tick, focus.Name)
		}

		if focus.Screen.View == nil {
			return node, fmt.Errorf(errf_view, focus.Name)
		}

		for c := range focus.Children() {
			if visited.Has(c.Id()) {
				return node, fmt.Errorf(errf_cycle, c.Name)
			}

			pending = append(pending, *c)
		}

		cursor++
	}

	return node, nil
}
