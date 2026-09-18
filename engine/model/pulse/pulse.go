package pulse

import (
	"time"

	assert "github.com/Rafael24595/go-assert/assert/runtime"
)

// Pulse wraps a time.Ticker to allow dynamic enabling, disabling, and safe cleanup of tick notifications.
type Pulse struct {
	chn    <-chan time.Time
	tkr    *time.Ticker
	closed bool
}

// New initializes a new Pulse instance with the specified ticker interval.
func New(duration time.Duration) *Pulse {
	return &Pulse{
		chn:    nil,
		tkr:    time.NewTicker(duration),
		closed: false,
	}
}

// Listen returns the channel where time ticks are received, or nil if the pulse is disabled.
// Triggers an assertion if called on a closed Pulse.
func (p *Pulse) Listen() <-chan time.Time {
	assert.False(p.closed, "cannot listen a closed pulse")

	return p.chn
}

// Enable activates tick channel output.
// Triggers an assertion if called on a closed Pulse.
func (p *Pulse) Enable() *Pulse {
	if p.closed {
		assert.Unreachable("closed pulse cannot be modified")
		return p
	}

	p.chn = p.tkr.C
	return p
}

// Disable suppresses tick channel output without stopping the internal ticker.
// Triggers an assertion if called on a closed Pulse.
func (p *Pulse) Disable() *Pulse {
	if p.closed {
		assert.Unreachable("closed pulse cannot be modified")
		return p
	}

	p.chn = nil
	return p
}

// Exit stops the internal ticker, detaches the channel, and marks the pulse as closed.
// Triggers an assertion if called on an already closed Pulse.
func (p *Pulse) Exit() *Pulse {
	if p.closed {
		assert.Unreachable("the pulse is already closed")
		return p
	}

	p.Disable()
	p.tkr.Stop()

	p.closed = true

	return p
}
