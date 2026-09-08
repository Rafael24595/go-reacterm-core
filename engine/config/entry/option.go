package entry

import (
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Option defines a functional configuration parameter for an Entry.
type Option func(*Entry)

func defaultEntry(node screen.Node) Entry {
	return Entry{
		Node:       node,
		Selectable: false,
		Opts:       make([]layer.Option[winsize.Rows], 0),
	}
}

// WithLayout assigns layout options for row sizing constraints.
func WithLayout(opts ...layer.Option[winsize.Rows]) Option {
	return func(cfg *Entry) {
		cfg.Opts = opts
	}
}

// Selectable marks the entry as focusable/interactable within the screen layout.
func Selectable() Option {
	return func(cfg *Entry) {
		cfg.Selectable = true
	}
}
