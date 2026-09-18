package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

const blinkMS = 750

// TextCursor manages caret position, text selection, and blink animation timing.
type TextCursor struct {
	clock  clock.Clock
	blink  bool
	status bool
	time   int64
	caret  offset.Offset
	anchor offset.Offset
}

// NewTextCursor initializes a new TextCursor with optional clock dependency.
func NewTextCursor(blink bool) *TextCursor {
	return NewTextCursorWithClock(blink, clock.UnixMilliClock)
}

// NewTextCursorWithClock initializes a new TextCursor with a custom clock provider.
func NewTextCursorWithClock(blink bool, clock clock.Clock) *TextCursor {
	return &TextCursor{
		clock:  clock,
		blink:  blink,
		status: true,
	}
}

// IsBlinking reports whether the cursor blinking feature is enabled.
func (c *TextCursor) IsBlinking() bool {
	return c.blink
}

// EnableBlinking enables blinking for the cursor.
func (c *TextCursor) EnableBlinking() *TextCursor {
	c.blink = true
	return c
}

// DisableBlinking disables blinking for the cursor.
func (c *TextCursor) DisableBlinking() *TextCursor {
	c.blink = false
	return c
}

// Caret returns the current caret offset.
func (c *TextCursor) Caret() offset.Offset {
	return c.caret
}

// Anchor returns the current selection anchor offset.
func (c *TextCursor) Anchor() offset.Offset {
	return c.anchor
}

// SelectStart returns the lower boundary offset of the current selection.
func (c *TextCursor) SelectStart() offset.Offset {
	if c.anchor < c.caret {
		return c.anchor
	}
	return c.caret
}

// SelectEnd returns the upper boundary offset of the current selection.
func (c *TextCursor) SelectEnd() offset.Offset {
	if c.anchor < c.caret {
		return c.caret
	}
	return c.anchor
}

// MoveCaretToStart moves the caret to the beginning of the buffer and resets blink timer.
func (c *TextCursor) MoveCaretToStart(buff []rune) *TextCursor {
	return c.MoveCaretToStartSilent(buff).Tick()
}

// MoveCaretToStartSilent moves the caret to the beginning of the buffer without resetting blink.
func (c *TextCursor) MoveCaretToStartSilent(buff []rune) *TextCursor {
	return c.MoveCaretTo(buff, 0)
}

// MoveCaretToEnd moves the caret to the end of the buffer and resets blink timer.
func (c *TextCursor) MoveCaretToEnd(buff []rune) *TextCursor {
	return c.MoveCaretToEndSilent(buff).Tick()
}

// MoveCaretToEndSilent moves the caret to the end of the buffer without resetting blink.
func (c *TextCursor) MoveCaretToEndSilent(buff []rune) *TextCursor {
	return c.MoveCaretTo(buff, offset.Offset(len(buff)))
}

// MoveCaretTo sets caret and anchor to a target position and resets blink timer.
func (c *TextCursor) MoveCaretTo(buff []rune, caret offset.Offset) *TextCursor {
	return c.MoveCaretSilent(buff, caret).Tick()
}

// MoveCaretSilent sets caret and anchor to a target position without resetting blink.
func (c *TextCursor) MoveCaretSilent(buff []rune, caret offset.Offset) *TextCursor {
	min := offset.Offset(1)
	len := offset.Offset(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = c.caret

	return c
}

// SelectAll selects the entire buffer and resets blink timer.
func (c *TextCursor) SelectAll(buff []rune) *TextCursor {
	return c.SelectAllSilent(buff).Tick()
}

// SelectAllSilent selects the entire buffer without resetting blink.
func (c *TextCursor) SelectAllSilent(buff []rune) *TextCursor {
	return c.SelectRangeSilent(buff, 0, offset.Offset(len(buff)))
}

// SelectRange updates caret and anchor boundaries and resets blink timer.
func (c *TextCursor) SelectRange(buff []rune, caret, anchor offset.Offset) *TextCursor {
	return c.SelectRangeSilent(buff, caret, anchor).Tick()
}

// SelectRangeSilent updates caret and anchor boundaries without resetting blink.
func (c *TextCursor) SelectRangeSilent(buff []rune, caret, anchor offset.Offset) *TextCursor {
	min := offset.Offset(1)
	len := offset.Offset(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = math.Clamp(anchor, min, len)

	return c
}

// Tick resets the blink cycle to visible status.
func (c *TextCursor) Tick() *TextCursor {
	c.status = true
	c.time = c.clock()
	return c
}

// BlinkStyle evaluates and returns the visual style state based on current time and blink state.
func (c *TextCursor) BlinkStyle() atom.Atom {
	if !c.blink || c.caret != c.anchor {
		return atom.Select
	}

	styl := atom.None
	if c.status {
		styl = atom.Select
	}

	now := c.clock()
	if now-c.time >= blinkMS {
		c.time = now
		c.status = !c.status
	}

	return styl
}
