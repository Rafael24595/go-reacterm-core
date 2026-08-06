package wrap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"
)

type Option func(*Wrapper)

func DefaultWrapper() Wrapper {
	return Wrapper{
		processors: []processor.Line{
			processor.LineFeed,
		},
		splitter: splitter.SplitLine,
	}
}

func WithProcessors(processors ...processor.Line) Option {
	return func(cfg *Wrapper) {
		cfg.processors = append(
			cfg.processors, processors...,
		)
	}
}

func WithSplitter(splitter splitter.Line) Option {
	return func(cfg *Wrapper) {
		cfg.splitter = splitter
	}
}
