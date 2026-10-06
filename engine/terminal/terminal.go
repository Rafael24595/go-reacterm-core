package terminal

import (
	"github.com/Rafael24595/go-reacterm-core/engine/model/key"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

// Terminal defines the operational contract and callbacks required to interact
// with an underlying terminal device or emulator session.
type Terminal struct {
	// OnStart is invoked when the terminal lifecycle begins.
	OnStart      func() error
	// OnClose is invoked when the terminal session terminates or cleans up resources.
	OnClose      func() error
	// ResizeEvents returns a read-only channel emitting terminal window size changes.
	ResizeEvents func() <-chan winsize.Winsize
	// KeyEvents returns a read-only channel emitting keyboard input events.
	KeyEvents    func() <-chan key.Key
	// Size queries and returns the current terminal window dimensions.
	Size         func() (winsize.Winsize, error)
	// Clear erases the terminal screen buffer content.
	Clear        func() error
	// Write renders individual string segments to the terminal output buffer.
	Write        func(...string) error
	// WriteAll renders a full raw string payload directly to the terminal output.
	WriteAll     func(string) error
	// Flush flushes any pending buffered output data directly to the active screen.
	Flush        func() error
}
