package textarea

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestTextArea_UnitBasicSuite(t *testing.T) {
	unit := New([]rune{}, input.NewTextCursor(false)).ToUnit()
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestTextAreaUnit_ResolveFrags(t *testing.T) {
	text := []rune("Hello World")

	cursor := input.NewTextCursor(false).
		SelectAll(text)

	got := New(text, cursor).
		resolveFrags()

	assert.Size(t, 1, got)
	assert.Equal(t, string(text), text_test.FragsToString(got))
}

func TestTextAreaUnit_ResolveFrags_WithEmptyBuffer(t *testing.T) {
	t.Run("empty buffer without placeholder generates default caret frags", func(t *testing.T) {
		text := []rune{}

		cursor := input.NewTextCursor(false).
			SelectAll(text)

		got := New(text, cursor).
			resolveFrags()

		assert.Size(t, 1, got)
		assert.Equal(t, marker.DefaultPaddingText, text_test.FragsToString(got))
	})

	t.Run("uses placeholder when buffer is empty and placeholder is present", func(t *testing.T) {
		text := []rune{}
		placeholder := "Type here..."

		cursor := input.NewTextCursor(false).
			SelectAll(text)

		got := New(text, cursor).
			SetPlaceholder(placeholder).
			resolveFrags()

		assert.Size(t, 1, got)
		assert.True(t, got[0].Atom().HasAny(atom.Dim))
		assert.Equal(t, placeholder, text_test.FragsToString(got))
	})

	t.Run("ignores placeholder when buffer contains text", func(t *testing.T) {
		text := []rune("golang")
		placeholder := "Type here..."

		cursor := input.NewTextCursor(false).
			SelectAll(text)

		got := New(text, cursor).
			SetPlaceholder(placeholder).
			resolveFrags()

		assert.Size(t, 1, got)
		assert.False(t, got[0].Atom().HasAny(atom.Dim))
		assert.Equal(t, string(text), text_test.FragsToString(got))
	})
}
