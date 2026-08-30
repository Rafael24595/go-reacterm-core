package behavior

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
)

func TestNewContext(t *testing.T) {
	target := Target{Name: "settings"}

	dummyNext := func(u state.UIState) {}
	ctx := NewContext[screen.BootFunc](target, dummyNext)

	assert.Equal(t, "settings", ctx.Target.Name)
	assert.NotNil(t, ctx.Next)
}
