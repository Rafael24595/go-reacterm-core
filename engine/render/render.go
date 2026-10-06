package render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// StringProcessor defines a function strategy that processes line slice data
// alongside container window dimensions into a single formatted output string.
type StringProcessor func([]line.Line, winsize.Winsize) string

// LinesProcessor defines a function strategy that processes line slice data
// alongside container window dimensions into a slice of formatted output strings.
type LinesProcessor func([]line.Line, winsize.Winsize) []string

// Render encapsulates the primary processing pipeline responsible for rendering
// text lines into formatted string output.
type Render struct {
	// Processor is the core function that transforms line slice data and window dimensions
	Processor StringProcessor
}

// RenderBuilder provides a fluent interface for configuring and constructing
// a Render instance.
type RenderBuilder struct {
	render StringProcessor
}

// NewBuilder constructs and initializes a new RenderBuilder with the specified
// StringProcessor function.
func NewBuilder(processor StringProcessor) *RenderBuilder {
	return &RenderBuilder{
		render: processor,
	}
}

// ToRender constructs and returns a Render value populated with the configured processor.
func (b *RenderBuilder) ToRender() Render {
	return Render{
		Processor: b.render,
	}
}
