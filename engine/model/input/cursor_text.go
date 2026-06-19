package input

import (
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/model/offset"
	"github.com/Rafael24595/go-reacterm-core/engine/platform/clock"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
)

const blink_ms = 750

type TextCursor struct {
	clock  clock.Clock
	blink  bool
	status bool
	time   int64
	caret  offset.Offset
	anchor offset.Offset
}

func NewTextCursor(blink bool) *TextCursor {
	return &TextCursor{
		clock:  clock.UnixMilliClock,
		blink:  blink,
		status: true,
		time:   0,
		caret:  0,
		anchor: 0,
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

func (c *TextCursor) MoveCaretTo(buff []rune, caret offset.Offset) {
	c.MoveCaretWithoutTick(buff, caret)
	c.Tick()
}

func (c *TextCursor) MoveCaretWithoutTick(buff []rune, caret offset.Offset) {
	min := offset.Offset(1)
	len := offset.Offset(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = c.caret
}

func (c *TextCursor) MoveSelectTo(buff []rune, caret, anchor offset.Offset) {
	c.MoveSelectWithoutTick(buff, caret, anchor)
	c.Tick()
}

func (c *TextCursor) MoveSelectWithoutTick(buff []rune, caret, anchor offset.Offset) {
	min := offset.Offset(1)
	len := offset.Offset(len(buff))

	if len == 0 {
		min = 0
	}

	c.caret = math.Clamp(caret, min, len)
	c.anchor = math.Clamp(anchor, min, len)
}

func (c *TextCursor) Tick() {
	c.status = true
	c.time = c.clock()
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
	if now-c.time >= blink_ms {
		c.time = now
		c.status = !c.status
	}

	return styl
}
