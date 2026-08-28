package pass

import (
	"errors"
	"fmt"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

var (
	// ErrMissingName indicates that a Node is missing a name.
	ErrMissingName = errors.New("screen: name is required")
	// ErrMissingBoot indicates that a Node is missing a Boot function.
	ErrNilBoot     = errors.New("screen: Boot function is nil")
	// ErrMissingKeys indicates that a Node is missing a Keys function.
	ErrNilKeys     = errors.New("screen: Keys function is nil")
	// ErrMissingTick indicates that a Node is missing a Tick function.
	ErrNilTick     = errors.New("screen: Tick function is nil")
	// ErrMissingView indicates that a Node is missing a View function.
	ErrNilView     = errors.New("screen: View function is nil")
	// ErrCycleFound indicates that a cycle was detected in the node tree.
	ErrCycleFound  = errors.New("screen: cycle detected in node tree")
)

// ValidateStructure is a compiler Pass that traverses the node hierarchy using breadth-first search (BFS).
// It verifies that each node in the tree has a valid name and fully initialized screen lifecycle functions
// (Boot, Keys, Tick, View), while ensuring no cyclic dependencies exist between parent and child nodes.
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
