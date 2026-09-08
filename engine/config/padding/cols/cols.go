package cols

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// FragProvider computes a column rendering fragment based on available width and input lines.
type FragProvider func(winsize.Cols, ...line.Line) frag.Frag

// Option defines a functional parameter to configure column behavior.
type Option func(*Config)

// Config holds horizontal position alignment and fragment generation strategies.
type Config struct {
	// Position specifies the horizontal alignment of the columns within the layout.
	Position style.HorizontalPosition
	// Provider generates a text fragment based on the current window width and existing lines.
	Provider FragProvider
}

// ResolveConfig creates a Config initialized with default settings and applies all provided Options.
func ResolveConfig(opts ...Option) Config {
	cfg := defaultColsConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func defaultColsConfig() Config {
	return Config{
		Position: style.Left,
		Provider: func(winsize.Cols, ...line.Line) frag.Frag {
			return frag.FromString(marker.DefaultPaddingText)
		},
	}
}

// WithPosition overrides the default horizontal alignment position.
func WithPosition(position style.HorizontalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

// WithText sets a static text value as the column fragment provider.
func WithText(txt string) Option {
	return WithFrag(
		frag.FromString(txt),
	)
}

// WithFrag configures a custom static frag.Frag instance as the column provider.
func WithFrag(frg frag.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(winsize.Cols, ...line.Line) frag.Frag {
			return frg
		}
	}
}
