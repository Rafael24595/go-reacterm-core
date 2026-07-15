package render

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type Processor func([]line.Line, winsize.Winsize) string
type RawProcessor func([]line.Line, winsize.Winsize) []string

type Render struct {
	Processor Processor
}

type RenderBuilder struct {
	render Processor
}

func NewBuilder(processor Processor) *RenderBuilder {
	return &RenderBuilder{
		render: processor,
	}
}

func (b *RenderBuilder) ToRender() Render {
	return Render{
		Processor: b.render,
	}
}
