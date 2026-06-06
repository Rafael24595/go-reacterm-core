package input

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestNormalizeMenuOptions_Success(t *testing.T) {
	input := []MenuOption{
		{Id: "lang"},
		{Id: "framework"},
		{Id: "tool"},
	}

	got := NormalizeMenuOptions(input...)

	assert.DeepEqual(t, input, got)
}

func TestNormalizeMenuOptions_TriggersAssertOnDuplicate(t *testing.T) {
	input := []MenuOption{
		{Id: "lang"},
		{Id: "lang"},
	}

	assert.Panic(t, func() {
		_ = NormalizeMenuOptions(input...)
	})
}
