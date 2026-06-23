package form

import (
	"testing"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestDefaultReadBindingsCoverAllCommands(t *testing.T) {
	screen_test.Helper_BindingsCover(
		t, defaultReadBindings, CommandsRead,
	)
}

func TestDefaultWriteBindingsCoverAllCommands(t *testing.T) {
	screen_test.Helper_BindingsCover(
		t, defaultWriteBindings, CommandsWrite,
	)
}
