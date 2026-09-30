package wrap

import (
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/processor"
	"github.com/Rafael24595/go-reacterm-core/engine/render/wrap/splitter"
)

// Option defines a functional configuration parameter for customizing a Wrapper instance.
type Option func(*Wrapper)

// DefaultWrapper returns a Wrapper initialized with default line feed processing and standard line splitting.
func DefaultWrapper() Wrapper {
	return Wrapper{
		processors: []processor.Line{
			processor.LineFeed,
		},
		splitter: splitter.SplitLine,
	}
}

// WithProcessors appends additional line processing functions to the Wrapper pipeline.
func WithProcessors(processors ...processor.Line) Option {
	return func(cfg *Wrapper) {
		cfg.processors = append(
			cfg.processors, processors...,
		)
	}
}

// WithSplitter configures a custom line splitting strategy for the Wrapper.
func WithSplitter(splitter splitter.Line) Option {
	return func(cfg *Wrapper) {
		cfg.splitter = splitter
	}
}
