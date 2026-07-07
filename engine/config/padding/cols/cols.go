package cols

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type FragProvider func(winsize.Cols, ...text.Line) text.Frag

type Option func(*Config)

type Config struct {
	Position style.HorizontalPosition
	Provider FragProvider
}

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
		Provider: func(_ winsize.Cols, _ ...text.Line) text.Frag {
			return *text.NewFrag(marker.DefaultPaddingText)
		},
	}
}

func WithPosition(position style.HorizontalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithText(frag string) Option {
	return WithFrag(
		*text.NewFrag(frag),
	)
}

func WithFrag(frag text.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Cols, _ ...text.Line) text.Frag {
			return frag
		}
	}
}
