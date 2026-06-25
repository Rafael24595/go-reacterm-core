package gutter

import "github.com/Rafael24595/go-reacterm-core/engine/render/marker"

const (
	DefaultLeft  = marker.U258C_Text + marker.DefaultPaddingText
	DefaultEmpty = marker.DefaultPaddingText + marker.DefaultPaddingText
)

type Option func(*meta)

type meta struct {
	left  string
	right string
}

func defaultMeta() meta {
	return meta{
		left:  DefaultLeft,
		right: "",
	}
}

func WithGutter(left, right string) Option {
	return func(cfg *meta) {
		cfg.left = left
		cfg.right = right
	}
}

func WithLeftGutter(gutter string) Option {
	return func(cfg *meta) {
		cfg.left = gutter
	}
}

func WithRightGutter(gutter string) Option {
	return func(cfg *meta) {
		cfg.right = gutter
	}
}

func newMeta(opts ...Option) meta {
	if len(opts) == 0 {
		return defaultMeta()
	}

	cfg := meta{}
	for _, o := range opts {
		o(&cfg)
	}
	return cfg
}
