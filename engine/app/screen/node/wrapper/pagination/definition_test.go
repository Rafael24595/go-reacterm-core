package pagination

import (
	"testing"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestDefaultPagedBindingsCoverAllCommands(t *testing.T) {
	screen_test.Helper_BindingsCover(
		t, defaultBindings.get(action.KindPaged), Commands,
	)
}

func TestDefaultScrollBindingsCoverAllCommands(t *testing.T) {
	screen_test.Helper_BindingsCover(
		t, defaultBindings.get(action.KindScroll), Commands,
	)
}
