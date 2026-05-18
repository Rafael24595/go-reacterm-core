package rule

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestAddSpaceAfter_AddsSpace(t *testing.T) {
	text, ok := AppendSpaceAfter(
		[]rune{','},
		5,
		5,
		nil,
	)

	assert.True(t, ok)
	assert.Equal(t, ", ", string(text))
}

func TestAddSpaceAfter_IgnoresOtherRunes(t *testing.T) {
	text, ok := AppendSpaceAfter(
		[]rune{'a'},
		1,
		1,
		nil,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}

func TestWrapSelection_WrapsSelectionWithBrackets(t *testing.T) {
	buffer := []rune("hello")

	text, ok := WrapSelection(
		[]rune{'('},
		0,
		5,
		buffer,
	)

	assert.True(t, ok)
	assert.Equal(t, "(hello)", string(text))
}

func TestWrapSelection_DoesNothingIfRuneIsNotWrapper(t *testing.T) {
	buffer := []rune("hello")

	text, ok := WrapSelection(
		[]rune{'a'},
		1,
		4,
		buffer,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}
