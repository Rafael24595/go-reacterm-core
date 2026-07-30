package wrap

type Option func(*Wrapper)

func defaultWrapper() Wrapper {
	return Wrapper{
		processors: []Processor{
			LineFeedProcessor,
		},
		splitter: SplitLine,
	}
}

/*func defaultWrapper() Wrapper {
	return Wrapper{
		processors: []Processor{
			LineFeedProcessor,
			ChunkProcessor(chunk.DefaultChunk),
		},
		splitter: CacheLineWords(
			NewFragCache(),
		),
	}
}*/

func WithProcessors(processors ...Processor) Option {
	return func(cfg *Wrapper) {
		cfg.processors = append(
			cfg.processors, processors...,
		)
	}
}

func WithSplitter(splitter LineSplitter) Option {
	return func(cfg *Wrapper) {
		cfg.splitter = splitter
	}
}
