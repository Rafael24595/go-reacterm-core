package format

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestTextFromString(t *testing.T) {
	txt := TextFromString("Hello")

	assert.Equal(t, "Hello", txt.Data)
	assert.Equal(t, winsize.Cols(5), txt.Size)
	assert.False(t, txt.IsEmpty())
}

func TestTextFromAny(t *testing.T) {
	txt := TextFromAny(12345)

	assert.Equal(t, "12345", txt.Data)
	assert.Equal(t, winsize.Cols(5), txt.Size)
	assert.Equal(t, "12345", txt.Data)
}

func TestEmptyText(t *testing.T) {
	txt := EmptyText()

	assert.Equal(t, "", txt.Data)
	assert.Equal(t, winsize.Cols(0), txt.Size)
	assert.True(t, txt.IsEmpty())
}
