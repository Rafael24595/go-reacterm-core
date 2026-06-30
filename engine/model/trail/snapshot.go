package trail

import "github.com/Rafael24595/go-reacterm-core/engine/app/screen"

type Snapshot struct {
	Previous []screen.Node
	Current  screen.Node
	Next     []screen.Node
}

func (s Snapshot) ToSlice() []screen.Node {
	items := append(s.Previous, s.Current)
	return append(items, s.Next...)
}
