package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

func TestNormalizeMenuOptions_Success(t *testing.T) {
	frg := frag.Empty()

	input := []MenuOption{
		NewMenuOption("lang", frg),
		NewMenuOption("framework", frg),
		NewMenuOption("tool", frg),
	}

	got := NormalizeMenuOptions(input...)

	assert.DeepEqual(t, input, got)
}

func TestNormalizeMenuOptions_TriggersAssertOnDuplicate(t *testing.T) {
	frg := frag.Empty()

	input := []MenuOption{
		NewMenuOption("lang", frg),
		NewMenuOption("lang", frg),
	}

	assert.Panic(t, func() {
		_ = NormalizeMenuOptions(input...)
	})
}
