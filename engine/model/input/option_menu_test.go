package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestMenuOption_InitializationAndGetters(t *testing.T) {
	frg := frag.FromString("Settings")
	opt := NewMenuOption("opt-1", frg)

	assert.Equal(t, "opt-1", opt.Id())

	assert.Equal(t,
		text_test.FragsToString(frg),
		text_test.FragsToString(opt.Label()),
	)
}

func TestMenuOption_WithHandlerAndExec(t *testing.T) {
	frg := frag.Empty()
	opt := NewMenuOption("opt-1", frg)

	node, ok := opt.Exec()
	assert.False(t, ok)
	assert.DeepEqual(t, screen.Node{}, node)

	optSame := opt.WithHandler(nil)
	_, ok = optSame.Exec()
	assert.False(t, ok)

	expectedNode := screen.Node{
		Name: "mock",
	}

	optWithHandler := opt.WithHandler(func() screen.Node {
		return expectedNode
	})

	node, ok = optWithHandler.Exec()
	assert.True(t, ok)
	assert.DeepEqual(t, expectedNode, node)
}

func TestExtractMenuOptionLabels(t *testing.T) {
	f1 := frag.FromString("Option 1")
	f2 := frag.FromString("Option 2")

	opts := []MenuOption{
		NewMenuOption("1", f1),
		NewMenuOption("2", f2),
	}

	labels := ExtractMenuOptionLabels(opts...)

	assert.Size(t, 2, labels)

	assert.Equal(t,
		text_test.FragsToString(f1),
		text_test.FragsToString(labels[0]),
	)

	assert.Equal(t,
		text_test.FragsToString(f2),
		text_test.FragsToString(labels[1]),
	)
}

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
