package behavior

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

type Target struct {
	Name string
	Tags set.Set[string]
}

func TargetOf(node screen.Node) Target {
	return Target{
		Name: node.Name,
		Tags: node.Tags,
	}
}
