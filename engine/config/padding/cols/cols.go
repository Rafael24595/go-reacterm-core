package cols

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type FragProvider func(winsize.Cols, ...text.Line) frag.Frag

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
		Provider: func(_ winsize.Cols, _ ...text.Line) frag.Frag {
			return *frag.New(marker.DefaultPaddingText)
		},
	}
}

func WithPosition(position style.HorizontalPosition) Option {
	return func(cfg *Config) {
		cfg.Position = position
	}
}

func WithText(txt string) Option {
	return WithFrag(
		*frag.New(txt),
	)
}

func WithFrag(frg frag.Frag) Option {
	return func(cfg *Config) {
		cfg.Provider = func(_ winsize.Cols, _ ...text.Line) frag.Frag {
			return frg
		}
	}
}
