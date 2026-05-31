package form

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
)

func TestPointer_Bitmasks(t *testing.T) {
	tests := []struct {
		name     string
		mask     Pointer
		checkAny []Pointer
		wantAny  bool
		wantNone bool
	}{
		{
			name:     "Prompt active has Prompt",
			mask:     PointerPrompt,
			checkAny: []Pointer{PointerPrompt},
			wantAny:  true,
			wantNone: false,
		},
		{
			name:     "Prompt active does not have Gutter",
			mask:     PointerPrompt,
			checkAny: []Pointer{PointerGutter},
			wantAny:  false,
			wantNone: true,
		},
		{
			name:     "Combo active has any of them",
			mask:     PointerPrompt | PointerGutter,
			checkAny: []Pointer{PointerGutter},
			wantAny:  true,
			wantNone: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantAny, tt.mask.HasAny(tt.checkAny...))
			assert.Equal(t, tt.wantNone, tt.mask.HasNone(tt.checkAny...))
		})
	}
}

func TestPointer_Navigation(t *testing.T) {
	t.Run("findPointer bounds checking", func(t *testing.T) {
		assert.Equal(t, PointerGutter, FindPointer(0))
		assert.Equal(t, PointerPrompt, FindPointer(1))

		assert.Equal(t, PointerGutter, FindPointer(3))
		assert.Equal(t, PointerGutter, FindPointer(255))
	})

	t.Run("nextPointer cycling logic", func(t *testing.T) {
		assert.Equal(t, uint8(1), NextPointer(0))
		assert.Equal(t, uint8(2), NextPointer(1))
		assert.Equal(t, uint8(0), NextPointer(2))
	})
}
