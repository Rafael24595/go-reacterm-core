package wrap

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/layout"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"
)

func dummyProcessor(order bool, lne line.Line) []line.Line {
	return []line.Line{lne}
}

func TestDefaultWrapper(t *testing.T) {
	w := DefaultWrapper()

	assert.Size(t, 1, w.processors)
	assert.NotNil(t, w.splitter)
}

func TestWithProcessors(t *testing.T) {
	w := NewWrapper(
		splitter.SplitLine,
		WithProcessors(dummyProcessor),
	)

	assert.Size(t, 1, w.processors)
}

func TestWithSplitter(t *testing.T) {
	customSplitterCalled := false

	customSplitter := func(lne line.Line) ([]layout.Word, []layout.Frag) {
		customSplitterCalled = true
		return splitter.SplitLine(lne)
	}

	w := NewWrapper(
		splitter.SplitLine,
		WithSplitter(customSplitter),
	)

	_ = w.normalize(false, line.FromString("test"))

	assert.True(t, customSplitterCalled)
}

func TestFromWrapper(t *testing.T) {
	base := DefaultWrapper()

	derived := FromWrapper(
		base, WithProcessors(dummyProcessor),
	)

	assert.Size(t, 1, base.processors)
	assert.Size(t, 2, derived.processors)
}
