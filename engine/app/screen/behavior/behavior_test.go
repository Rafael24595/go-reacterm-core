package behavior

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestApply(t *testing.T) {
	initialNode := screen.Node{Name: "initial"}

	renameBehavior := func(n screen.Node) screen.Node {
		n.Name = n.Name + "_mod_01"
		return n
	}

	tagBehavior := func(n screen.Node) screen.Node {
		if n.Tags == nil {
			n.Tags = set.New[string]()
		}
		n.Tags.Add("processed")
		return n
	}

	resultNode := Apply(initialNode, renameBehavior, tagBehavior)

	assert.Equal(t, "initial_mod_01", resultNode.Name)
	
	assert.NotNil(t, resultNode.Tags)
	assert.True(t, resultNode.Tags.Has("processed"))
}

func TestApply_EmptyBehaviors(t *testing.T) {
	initialNode := screen.Node{Name: "unmodified"}

	resultNode := Apply(initialNode)

	assert.Equal(t, "unmodified", resultNode.Name)
}
