package entry

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Entry represents a single renderable node within a screen layout layer.
type Entry struct {
	// Node is the renderable screen node associated with this entry.
	Node       screen.Node
	// Selectable indicates whether the entry can be selected or focused.
	Selectable bool
	// Opts holds additional configuration options for the entry's rendering behavior.
	Opts       []layer.Option[winsize.Rows]
}

// New creates and configures a new Entry with the given options.
// Panics if the provided node is a zero value.
func New(node screen.Node, opts ...Option) Entry {
	cfg := defaultEntry(node)
	for _, opt := range opts {
		opt(&cfg)
	}

	assert.LazyFalse(func() bool {
		return screen.IsZeroNode(cfg.Node)
	}, "unit is not defined")

	return cfg
}
