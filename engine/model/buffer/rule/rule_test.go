package rule

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestAutoSpace_AddsSpace(t *testing.T) {
	text, ok := AutoSpace(
		[]rune{','},
		5,
		5,
		nil,
	)

	assert.True(t, ok)
	assert.Equal(t, ", ", string(text))
}

func TestAutoSpace_IgnoresOtherRunes(t *testing.T) {
	text, ok := AutoSpace(
		[]rune{'a'},
		1,
		1,
		nil,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}

func TestAutoWrap_WrapsSelectionWithBrackets(t *testing.T) {
	buffer := []rune("hello")

	text, ok := AutoWrap(
		[]rune{'('},
		0,
		5,
		buffer,
	)

	assert.True(t, ok)
	assert.Equal(t, "(hello)", string(text))
}

func TestAutoWrap_DoesNothingIfRuneIsNotWrapper(t *testing.T) {
	buffer := []rune("hello")

	text, ok := AutoWrap(
		[]rune{'a'},
		1,
		4,
		buffer,
	)

	assert.False(t, ok)
	assert.Equal(t, "a", string(text))
}
