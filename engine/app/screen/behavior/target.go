package behavior

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

// Target holds metadata about a screen node used to match and target behavior interceptors.
type Target struct {
	// Name represents the identifier name of the targeted screen node.
	Name string
	// Tags holds a set of tags associated with the targeted screen node for category-based interception.
	Tags set.Set[string]
}

// TargetOf extracts metadata from a screen Node and returns a Target instance representing it.
func TargetOf(node screen.Node) Target {
	return Target{
		Name: node.Name,
		Tags: node.Tags,
	}
}
