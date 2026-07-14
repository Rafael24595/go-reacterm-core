package rows

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type FragProvider func(winsize.Winsize, ...text.Line) frag.Frag

type Option func(*Config)

type Config struct {
	Position style.VerticalPosition
	Provider FragProvider
}

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
		Provider: func(_ winsize.Winsize, _ ...text.Line) frag.Frag {
			return *frag.Empty()
		},
	}
}

func WithPosition(position style.VerticalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithFrag(frg frag.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Winsize, _ ...text.Line) frag.Frag {
			return frg
		}
	}
}

func WithFillFrag(txt ...string) Option {
	data := marker.DefaultPaddingText
	if len(txt) > 0 {
		data = txt[0]
	}

	return func(cfg *Config) {
		cfg.Provider = func(size winsize.Winsize, lines ...text.Line) frag.Frag {
			measure := text.MaxLineMeasure(size.Cols, lines...)
			return *frag.New(data).
				AddSpec(spec.ExtendRight(measure))
		}
	}
}
