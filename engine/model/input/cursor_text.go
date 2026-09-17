package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

const blinkMS = 750

type TextCursor struct {
	clock  clock.Clock
	blink  bool
	status bool
	time   int64
	caret  offset.Offset
	anchor offset.Offset
}

func NewTextCursor(blink bool) *TextCursor {
	return NewTextCursorWithClock(blink, clock.UnixMilliClock)
}

func NewTextCursorWithClock(blink bool, clock clock.Clock) *TextCursor {
	return &TextCursor{
		clock:  clock,
		blink:  blink,
		status: true,
	}
}

func (c *TextCursor) IsBlinking() bool {
	return c.blink
}

func (c *TextCursor) EnableBlinking() *TextCursor {
	c.blink = true
	return c
}

func (c *TextCursor) DisableBlinking() *TextCursor {
	c.blink = false
	return c
}

func (c *TextCursor) Caret() offset.Offset {
	return c.caret
}

func (c *TextCursor) Anchor() offset.Offset {
	return c.anchor
}

func (c *TextCursor) SelectStart() offset.Offset {
	if c.anchor < c.caret {
		return c.anchor
	}
	return c.caret
}

func (c *TextCursor) SelectEnd() offset.Offset {
	if c.anchor < c.caret {
		return c.caret
	}
	return c.anchor
}

func (c *TextCursor) MoveCaretToStart(buff []rune) *TextCursor {
	return c.MoveCaretToStartSilent(buff).Tick()
}

func (c *TextCursor) MoveCaretToStartSilent(buff []rune) *TextCursor {
	return c.MoveCaretTo(buff, 0)
}

func (c *TextCursor) MoveCaretToEnd(buff []rune) *TextCursor {
	return c.MoveCaretToEndSilent(buff).Tick()
}

func (c *TextCursor) MoveCaretToEndSilent(buff []rune) *TextCursor {
	return c.MoveCaretTo(buff, offset.Offset(len(buff)))
}

func (c *TextCursor) MoveCaretTo(buff []rune, caret offset.Offset) *TextCursor {
	return c.MoveCaretSilent(buff, caret).Tick()
}

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

func (c *TextCursor) SelectAll(buff []rune) *TextCursor {
	return c.SelectAllSilent(buff).Tick()
}

func (c *TextCursor) SelectAllSilent(buff []rune) *TextCursor {
	return c.SelectRangeSilent(buff, 0, offset.Offset(len(buff)))
}

func (c *TextCursor) SelectRange(buff []rune, caret, anchor offset.Offset) *TextCursor {
	return c.SelectRangeSilent(buff, caret, anchor).Tick()
}

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

func (c *TextCursor) Tick() *TextCursor {
	c.status = true
	c.time = c.clock()
	return c
}

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
