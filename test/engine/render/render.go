package render_test

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type MockRender struct {
	Processor render.Processor
}

func FromProcessor(processor render.Processor) render.Render {
	return MockRender{
		Processor: processor,
	}.ToRender()
}

func DiscardRender() render.Render {
	return MockRender{}.ToRender()
}

func (m MockRender) ToRender() render.Render {
	if m.Processor == nil {
		m.Processor = func([]line.Line, winsize.Winsize) string {
			return ""
		}
	}
	return render.Render{
		Processor: m.Processor,
	}
}
