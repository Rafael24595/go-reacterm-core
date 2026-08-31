package draw

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/work"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

// State maintains the frame buffer, asynchronous work tracking, and pagination status for a render pass.
type State struct {
	Buffer []line.Line
	Work   *work.Tracker
	Cursor uint16
	Page   uint
	Focus  bool
}

// NewState creates a new render State with an optional initial buffer capacity based on window rows.
func NewState(size ...winsize.Rows) *State {
	buffSize := winsize.Rows(0)
	if len(size) > 0 {
		buffSize = size[0]
	}

	return &State{
		Buffer: make([]line.Line, buffSize),
		Work:   work.NewTracker(),
		Cursor: 0,
		Page:   0,
		Focus:  false,
	}
}

// IsPaginated reports whether pagination controls should be rendered for the current state.
func (s *State) IsPaginated() bool {
	return s.Page != 0 || s.Work.Unfinished()
}

// EnableFocus updates the focus state using a logical OR, ensuring focus stays true once set.
func (s *State) EnableFocus(focus bool) *State {
	s.Focus = s.Focus || focus
	return s
}

// WriteLine appends a line to the buffer at the current cursor index and advances the cursor.
func (s *State) WriteLine(lne line.Line) *State {
	if s.IsFull() {
		assert.Unreachable("buffer overflow")
		return s
	}

	s.Buffer[s.Cursor] = lne
	s.Cursor += 1
	return s
}

// IsFull checks if the current cursor position has reached the buffer capacity.
func (s *State) IsFull() bool {
	return s.Cursor == uint16(len(s.Buffer))
}

// Written returns a slice of all lines written to the buffer up to the current cursor position.
func (s *State) Written() []line.Line {
	return s.Buffer[:s.Cursor]
}

// Reset clears the buffer elements and resets internal counters for frame reuse.
func (s *State) Reset() {
	clear(s.Buffer)
	s.Cursor = 0
	s.Focus = false
}
