package buffer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestClipboard_NewEmpty(t *testing.T) {
	cb := NewClipboard()

	assert.Equal(t, uint(0), cb.Size())
	assert.Nil(t, cb.Read())
}

func TestClipboard_WriteAndRead(t *testing.T) {
	cb := NewClipboard()
	input := []rune("hello world")

	cb.Write(input)

	assert.Equal(t, uint(11), cb.Size())
	assert.Equal(t, "hello world", string(cb.Read()))
}

func TestClipboard_ReadDefensiveCopy(t *testing.T) {
	cb := NewClipboard()
	cb.Write([]rune("golang"))

	readBuffer := cb.Read()
	readBuffer[0] = 'Z'

	assert.Equal(t, "golang", string(cb.Read()))
}

func TestClipboard_WriteDefensiveCopy(t *testing.T) {
	cb := NewClipboard()
	input := []rune("original")

	cb.Write(input)
	input[0] = 'X'

	assert.Equal(t, "original", string(cb.Read()))
}

func TestClipboard_Clear(t *testing.T) {
	cb := NewClipboard()
	cb.Write([]rune("data to clear"))

	cb.Clear()

	assert.Equal(t, uint(0), cb.Size())
	assert.Nil(t, cb.Read())
}

func TestClipboard_WriteFluentChain(t *testing.T) {
	cb := NewClipboard()

	res := cb.Write([]rune("first")).
		Write([]rune("second"))

	assert.Equal(t, cb, res)
	assert.Equal(t, "second", string(cb.Read()))
}
