package pass

import (
	"errors"
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

var (
	ErrMissingName = errors.New("screen: name is required")
	ErrNilBoot     = errors.New("screen: Boot function is nil")
	ErrNilKeys     = errors.New("screen: Keys function is nil")
	ErrNilTick     = errors.New("screen: Tick function is nil")
	ErrNilView     = errors.New("screen: View function is nil")
	ErrCycleFound  = errors.New("screen: cycle detected in node tree")
)

func ValidateStructure(node screen.Node) (screen.Node, error) {
	visited := set.New[string]()

	pending := []screen.Node{node}
	cursor := 0

	for cursor < len(pending) {
		focus := pending[cursor]
		visited.Add(focus.Id())

		if focus.Name == "" {
			return node, ErrMissingName
		}

		if focus.Screen.Boot == nil {
			return node, fmt.Errorf("screen %q: %w", focus.Name, ErrNilBoot)
		}

		if focus.Screen.Keys == nil {
			return node, fmt.Errorf("screen %q: %w", focus.Name, ErrNilKeys)
		}

		if focus.Screen.Tick == nil {
			return node, fmt.Errorf("screen %q: %w", focus.Name, ErrNilTick)
		}

		if focus.Screen.View == nil {
			return node, fmt.Errorf("screen %q: %w", focus.Name, ErrNilView)
		}

		for c := range focus.Children() {
			if visited.Has(c.Id()) {
				return node, fmt.Errorf("screen %q: %w", c.Name, ErrCycleFound)
			}

			pending = append(pending, *c)
		}

		cursor++
	}

	return node, nil
}
