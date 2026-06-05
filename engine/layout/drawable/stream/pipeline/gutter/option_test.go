package gutter

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestDefaultMeta(t *testing.T) {
	cfg := defaultMeta()

	assert.Equal(t, DefaultLeft, cfg.left)
	assert.Size(t, 0, cfg.right)
}

func TestWithGutter(t *testing.T) {
	cfg := newMeta(
		WithGutter(">", "<"),
	)

	assert.Equal(t, ">", cfg.left)
	assert.Equal(t, "<", cfg.right)
}

func TestWithLeftGutter(t *testing.T) {
	cfg := newMeta(
		WithLeftGutter(">"),
	)

	assert.Equal(t, ">", cfg.left)
	assert.Size(t, 0, cfg.right)
}

func TestWithRightGutter(t *testing.T) {
	cfg := newMeta(
		WithRightGutter("<"),
	)

	assert.Size(t, 0, cfg.left)
	assert.Equal(t, "<", cfg.right)
}
