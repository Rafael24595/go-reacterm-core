package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
)

// CheckOption represents an option entry with selection status and styled text fragment.
type CheckOption struct {
	id        string
	status    bool
	label     frag.Frag
	timestamp int64
}

// NewCheckOption initializes a new CheckOption with an ID and label fragment.
func NewCheckOption(id string, label frag.Frag) CheckOption {
	return CheckOption{
		id:    id,
		label: label,
	}
}

// Id returns the option unique identifier.
func (o *CheckOption) Id() string {
	return o.id
}

// Status reports whether the option is currently checked.
func (o *CheckOption) Status() bool {
	return o.status
}

// Label returns the text fragment label of the option.
func (o *CheckOption) Label() frag.Frag {
	return o.label
}

// Timestamp returns the unix milliseconds when the option was last checked.
func (o *CheckOption) Timestamp() int64 {
	return o.timestamp
}

// WithCheck returns a copy with updated status, setting timestamp via clock if checked.
func (o *CheckOption) WithCheck(status bool, clk clock.Clock) *CheckOption {
	o.status = status
	if o.status && clk != nil {
		o.timestamp = clk()
	}
	return o
}

// Check returns a copy set to checked (status = true) with the current timestamp.
func (o *CheckOption) Check(clk clock.Clock) *CheckOption {
	return o.WithCheck(true, clk)
}

// Uncheck returns a copy set to unchecked (status = false).
func (o *CheckOption) Uncheck() *CheckOption {
	o.status = false
	return o
}

// ExtractCheckOptionLabels extracts the underlying frag.Frag labels from a collection of CheckOptions.
func ExtractCheckOptionLabels(options ...CheckOption) []frag.Frag {
	lines := make([]frag.Frag, len(options))
	for i := range options {
		lines[i] = options[i].label
	}
	return lines
}
