package rows

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// FragProvider computes a rendering fragment based on window dimensions and existing text lines.
type FragProvider func(winsize.Winsize, ...line.Line) frag.Frag

// Option defines a functional parameter to configure row behavior.
type Option func(*Config)

// Config holds positioning and fragment rendering strategy for layout rows.
type Config struct {
	// Position specifies the vertical alignment of the rows within the layout.
	Position style.VerticalPosition
	// Provider generates a text fragment based on the current window size and existing lines.
	Provider FragProvider
}

// ResolveConfig creates a Config initialized with default settings and applies all provided Options.
func ResolveConfig(opts ...Option) Config {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

func defaultConfig() Config {
	return Config{
		Position: style.Top,
		Provider: func(winsize.Winsize, ...line.Line) frag.Frag {
			return frag.Empty()
		},
	}
}

// WithPosition overrides the default vertical alignment position.
func WithPosition(position style.VerticalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

// WithFrag sets a static fragment provider returning the given frag.Frag instance.
func WithFrag(frg frag.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(winsize.Winsize, ...line.Line) frag.Frag {
			return frg
		}
	}
}

// WithFillFrag configures a dynamic provider that pads row lines using the specified text pattern.
// Defaults to marker.DefaultPaddingText if omitted or empty.
func WithFillFrag(txt ...string) Option {
	data := marker.DefaultPaddingText
	if len(txt) > 0 && txt[0] != "" {
		data = txt[0]
	}

	return func(cfg *Config) {
		cfg.Provider = func(size winsize.Winsize, lines ...line.Line) frag.Frag {
			measure := line.MaxMeasure(size.Cols, lines...)
			return frag.TextSpec(data, spec.ExtendRight(measure))
		}
	}
}
