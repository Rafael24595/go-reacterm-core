package indexmenu

import (
	"testing"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestDefaultBindingsCoverAllCommands(t *testing.T) {
	screen_test.Helper_BindingsCover(t, defaultBindings, Commands)
}
