package render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type StringProcessor func([]line.Line, winsize.Winsize) string
type LinesProcessor func([]line.Line, winsize.Winsize) []string

type Render struct {
	Processor StringProcessor
}

type RenderBuilder struct {
	render StringProcessor
}

func NewBuilder(processor StringProcessor) *RenderBuilder {
	return &RenderBuilder{
		render: processor,
	}
}

func (b *RenderBuilder) ToRender() Render {
	return Render{
		Processor: b.render,
	}
}
