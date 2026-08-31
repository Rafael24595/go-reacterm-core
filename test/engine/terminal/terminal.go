package terminal_test

import (
	"strings"

	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/terminal"
)

type MockTerminal struct {
	Resize chan winsize.Winsize
	Keys   chan key.Key
	Size   winsize.Winsize
	Buffer strings.Builder
}

func DiscardTerminal() terminal.Terminal {
	return MockTerminal{}.ToTerminal(winsize.New(0, 0))
}

func (m MockTerminal) ToTerminal(size winsize.Winsize) terminal.Terminal {
	if m.Resize == nil {
		m.Resize = make(chan winsize.Winsize, 8)
	}

	if m.Keys == nil {
		m.Keys = make(chan key.Key, 64)
	}

	if size.Cols == 0 && size.Rows == 0 {
		m.Size = winsize.New(80, 200)
	}

	return terminal.Terminal{
		OnStart: func() error {
			return nil
		},
		OnClose: func() error {
			close(m.Resize)
			close(m.Keys)
			return nil
		},
		ResizeEvents: func() <-chan winsize.Winsize {
			return m.Resize
		},
		KeyEvents: func() <-chan key.Key {
			return m.Keys
		},
		Size: func() (winsize.Winsize, error) {
			return m.Size, nil
		},
		Clear: func() error {
			m.Buffer.Reset()
			return nil
		},
		Write: func(s ...string) error {
			for _, v := range s {
				m.Buffer.WriteString(v)
			}
			return nil
		},
		WriteAll: func(s string) error {
			m.Buffer.WriteString(s)
			return nil
		},
		Flush: func() error {
			return nil
		},
	}
}
