package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/test/support/mock"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestCheckOption_InitializationAndGetters(t *testing.T) {
	label := frag.FromString("Option 1")
	opt := NewCheckOption("opt-1", label)

	assert.Equal(t, "opt-1", opt.Id())
	assert.False(t, opt.Status())
	assert.Equal(t, 0, opt.Timestamp())

	assert.Equal(t,
		text_test.FragsToString(label),
		text_test.FragsToString(opt.label),
	)
}

func TestCheckOption_CheckAndUncheckSemantics(t *testing.T) {
	clk := &mock.TestClock{Time: 1000}

	label := frag.FromString("Option 1")
	opt := NewCheckOption("opt-1", label)

	checked := opt.Check(clk.Now)

	assert.True(t, checked.Status())
	assert.Equal(t, 1000, checked.Timestamp())

	clk.Advance(500)
	unchecked := checked.Uncheck()

	assert.False(t, unchecked.Status())
	assert.Equal(t, 1000, unchecked.Timestamp())
}

func TestCheckOption_WithCheckBehavior(t *testing.T) {
	clk := &mock.TestClock{Time: 2000}

	label := frag.FromString("Option 1")
	opt := NewCheckOption("opt-1", label)

	checked := opt.WithCheck(true, clk.Now)

	assert.True(t, checked.Status())
	assert.Equal(t, 2000, checked.Timestamp())

	clk.Advance(500)
	withFalse := checked.WithCheck(false, clk.Now)

	assert.False(t, withFalse.Status())
	assert.Equal(t, 2000, withFalse.Timestamp())
}

func TestCheckOption_NilClockGuard(t *testing.T) {
	label := frag.FromString("Option 1")
	opt := NewCheckOption("opt-1", label)

	checked := opt.Check(nil)

	assert.True(t, checked.Status())
	assert.Equal(t, 0, checked.Timestamp())
}

func TestCheckOption_ExtractLabels(t *testing.T) {
	f1 := frag.FromString("Item 1")
	f2 := frag.FromString("Item 2")
	f3 := frag.FromString("Item 3")

	opt1 := NewCheckOption("1", f1)
	opt2 := NewCheckOption("2", f2)
	opt3 := NewCheckOption("3", f3)

	labels := ExtractCheckOptionLabels(opt1, opt2, opt3)
	assert.Size(t, 3, labels)

	assert.Equal(t,
		text_test.FragsToString(f1),
		text_test.FragsToString(labels[0]),
	)

	assert.Equal(t,
		text_test.FragsToString(f2),
		text_test.FragsToString(labels[1]),
	)

	assert.Equal(t,
		text_test.FragsToString(f3),
		text_test.FragsToString(labels[2]),
	)

	emptyLabels := ExtractCheckOptionLabels()
	assert.Empty(t, emptyLabels)
}
