package rows

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type FragProvider func(winsize.Winsize, ...text.Line) text.Frag

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
		Provider: func(_ winsize.Winsize, _ ...text.Line) text.Frag {
			return *text.EmptyFrag()
		},
	}
}

func WithPosition(position style.VerticalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithFrag(frag text.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Winsize, _ ...text.Line) text.Frag {
			return frag
		}
	}
}

func WithFillFrag(frag ...string) Option {
	data := marker.DefaultPaddingText
	if len(frag) > 0 {
		data = frag[0]
	}

	return func(cfg *Config) {
		cfg.Provider = func(size winsize.Winsize, lines ...text.Line) text.Frag {
			measure := text.MaxLineMeasure(size.Cols, lines...)
			return *text.NewFrag(data).
				AddSpec(spec.ExtendRight(measure))
		}
	}
}
