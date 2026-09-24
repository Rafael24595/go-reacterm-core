package trail

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Snapshot struct {
	Previous []screen.Node
	Current  screen.Node
	Next     []screen.Node
}

func (s Snapshot) ToSlice() []screen.Node {
	total := len(s.Previous) + 1 + len(s.Next)
	result := make([]screen.Node, 0, total)

	result = append(result, s.Previous...)
	result = append(result, s.Current)
	result = append(result, s.Next...)

	return result
}
