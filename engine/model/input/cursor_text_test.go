package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/test/support/mock"
)

func TestCursor_SelectionLogic(t *testing.T) {
	c := NewTextCursor(true)
	buff := []rune("Golang")

	c.SelectAll(buff)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 6)

	c.SelectRange(buff, 3, 1)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 3)

	c.SelectRange(buff, 1, 3)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 3)

	c.MoveCaretToStart(buff)
	assert.Equal(t, c.Caret(), 1)

	c.MoveCaretToEnd(buff)
	assert.Equal(t, c.Caret(), 6)
}

func TestCursor_BlinkingLogic(t *testing.T) {
	clock := &mock.TestClock{Time: 0}

	c := NewTextCursor(true)
	c.clock = clock.Now

	clock.Advance(blinkMS + 1)

	assert.Equal(t, c.BlinkStyle(), atom.Select)

	clock.Advance(blinkMS + 1)

	assert.Equal(t, c.BlinkStyle(), atom.None)

	clock.Advance(blinkMS + 1)
	assert.Equal(t, c.BlinkStyle(), atom.Select)
}

func TestCursor_EmptyBufferAndClamping(t *testing.T) {
	c := NewTextCursor(true)

	c.MoveCaretTo([]rune(""), 5)
	assert.Equal(t, c.Caret(), 0)
	assert.Equal(t, c.Anchor(), 0)

	buff := []rune("Golang")

	c.MoveCaretTo(buff, 0)
	assert.Equal(t, c.Caret(), 1)

	c.MoveCaretTo(buff, 100)
	assert.Equal(t, c.Caret(), 6)

	c.SelectRange(buff, 0, 50)
	assert.Equal(t, c.SelectStart(), 1)
	assert.Equal(t, c.SelectEnd(), 6)
}

func TestCursor_BlinkOverrideConditions(t *testing.T) {
	clock := &mock.TestClock{Time: 0}
	c := NewTextCursorWithClock(true, clock.Now)

	buff := []rune("Golang")

	c.SelectRange(buff, 1, 4)
	clock.Advance(blinkMS + 1)

	assert.Equal(t, c.BlinkStyle(), atom.Select)

	c.MoveCaretTo(buff, 2).DisableBlinking()
	assert.False(t, c.IsBlinking())

	clock.Advance(blinkMS + 1)
	assert.Equal(t, c.BlinkStyle(), atom.Select)

	c.EnableBlinking()
	assert.True(t, c.IsBlinking())
}

func TestCursor_SilentOperations(t *testing.T) {
	clock := &mock.TestClock{Time: 0}
	c := NewTextCursorWithClock(true, clock.Now)
	
	buff := []rune("Golang")

	clock.Advance(blinkMS + 1)
	assert.Equal(t, c.BlinkStyle(), atom.Select)

	c.MoveCaretSilent(buff, 3)
	assert.Equal(t, c.BlinkStyle(), atom.None)

	c.MoveCaretTo(buff, 4)
	assert.Equal(t, c.BlinkStyle(), atom.Select)
}
