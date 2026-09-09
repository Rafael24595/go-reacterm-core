package format

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
)

func TestNewEllipsis(t *testing.T) {
	e := NewEllipsis(".", 3)

	assert.Equal(t, ".", e.Data)
	assert.Equal(t, winsize.Cols(3), e.Count)
	assert.Equal(t, winsize.Cols(3), e.measure())
	assert.Equal(t, "...", e.string())
}

func TestDefaultEllipsis(t *testing.T) {
	e := NewEllipsis(".", 3)

	assert.Equal(t, ".", e.Data)
	assert.Equal(t, winsize.Cols(3), e.Count)
	assert.Equal(t, "...", e.string())
}

func TestEllipsisWideRunes(t *testing.T) {
	e := NewEllipsis("…", 1)

	assert.Equal(t, winsize.Cols(1), e.measure())
	assert.Equal(t, "…", e.string())
}