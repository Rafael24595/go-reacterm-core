package draw

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/work"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

type State struct {
	Buffer []text.Line
	Work   *work.Tracker
	Cursor uint16
	Page   uint
	Focus  bool
}

func NewState(size ...winsize.Rows) *State {
	buffSize := winsize.Rows(0)
	if len(size) > 0 {
		buffSize = size[0]
	}

	return &State{
		Buffer: make([]text.Line, buffSize),
		Work:   work.NewTracker(),
		Cursor: 0,
		Page:   0,
		Focus:  false,
	}
}

func (s *State) ShowPagination() bool {
	return s.Page != 0 || s.Work.Unfinished()
}

func (s *State) MarkFocus(focus bool) *State {
	s.Focus = s.Focus || focus
	return s
}

func (s *State) SetAndNext(line text.Line) *State {
	if s.IsFull() {
		assert.Unreachable("buffer overflow")
		return s
	}

	s.Buffer[s.Cursor] = line
	s.Cursor += 1
	return s
}

func (s *State) IsFull() bool {
	return s.Cursor == uint16(len(s.Buffer))
}

func (s *State) Written() []text.Line {
	return s.Buffer[:s.Cursor]
}

func (s *State) Reset() {
	for i := range s.Buffer {
		s.Buffer[i] = text.Line{}
	}

	s.Cursor = 0
	s.Focus = false
}
