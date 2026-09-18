package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

type CheckOption struct {
	id        string
	status    bool
	label     frag.Frag
	timestamp int64
}

func NewCheckOption(id string, label frag.Frag) CheckOption {
	return CheckOption{
		id:    id,
		label: label,
	}
}

func (o *CheckOption) Id() string {
	return o.id
}

func (o *CheckOption) Status() bool {
	return o.status
}

func (o *CheckOption) Label() frag.Frag {
	return o.label
}

func (o *CheckOption) Timestamp() int64 {
	return o.timestamp
}

func (o *CheckOption) WithCheck(status bool, clk clock.Clock) *CheckOption {
	o.status = status
	if o.status && clk != nil {
		o.timestamp = clk()
	}
	return o
}

func (o *CheckOption) Check(clk clock.Clock) *CheckOption {
	return o.WithCheck(true, clk)
}

func (o *CheckOption) Uncheck() *CheckOption {
	o.status = false
	return o
}

func ExtractCheckOptionLabels(options ...CheckOption) []frag.Frag {
	lines := make([]frag.Frag, len(options))
	for i := range options {
		lines[i] = options[i].label
	}
	return lines
}
