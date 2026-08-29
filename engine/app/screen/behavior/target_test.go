package behavior

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/set"
)

func TestTargetOf(t *testing.T) {
	tags := set.New[string]()
	tags.Add("modal")
	tags.Add("focusable")

	node := screen.Node{
		Name: "main_screen",
		Tags: tags,
	}

	target := TargetOf(node)

	assert.Equal(t, "main_screen", target.Name)
	assert.True(t, target.Tags.Has("modal"))
	assert.True(t, target.Tags.Has("focusable"))
}
