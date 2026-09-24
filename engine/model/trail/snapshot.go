package trail

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

// Snapshot represents a point-in-time state of the navigation history trail.
type Snapshot struct {
	// Previous is the list of screens visited before the current screen, in order from oldest to most recent.
	Previous []screen.Node
	// Current is the screen currently being viewed.
	Current  screen.Node
	// Next is the list of screens that can be navigated to after the current screen, in order from most recent to oldest.
	Next     []screen.Node
}

// ToSlice flattens the history snapshot into a single ordered slice [Previous..., Current, Next...].
// It allocates a new slice to prevent mutating the underlying arrays.
func (s Snapshot) ToSlice() []screen.Node {
	total := len(s.Previous) + 1 + len(s.Next)
	result := make([]screen.Node, 0, total)

	result = append(result, s.Previous...)
	result = append(result, s.Current)
	result = append(result, s.Next...)

	return result
}
