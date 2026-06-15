package wrapper_render

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	
	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
)

func TestAtomsStyler(t *testing.T) {
	tests := []struct {
		name     string
		atom     atom.Atom
		input    string
		expected string
	}{
		{
			name:     "bold transformation",
			atom:     atom.Bold,
			input:    "Golang",
			expected: wrapper_ansi.Bold + "Golang" + wrapper_ansi.NormalWeight,
		},
		{
			name:     "select transformation",
			atom:     atom.Select,
			input:    "Ziglang",
			expected: wrapper_ansi.Reverse + "Ziglang" + wrapper_ansi.NoReverse,
		},
		{
			name:     "empty string",
			atom:     atom.Bold,
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := Atoms.Get(tt.atom)
			assert.True(t, ok)

			result := fn(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
