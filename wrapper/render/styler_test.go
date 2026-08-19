package wrapper_render

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	wrapper_ansi "github.com/Rafael24595/go-reacterm-core/wrapper/ansi"
)

func TestToBold(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "text",
			text: "Hello World",
			want: wrapper_ansi.Bold + "Hello World" + wrapper_ansi.NormalWeight,
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toBold(tt.text))
		})
	}
}

func TestToDim(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "text",
			text: "Hello World",
			want: wrapper_ansi.Dim + "Hello World" + wrapper_ansi.NormalWeight,
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toDim(tt.text))
		})
	}
}

func TestToSelect(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "text",
			text: "Hello World",
			want: wrapper_ansi.Reverse + "Hello World" + wrapper_ansi.NoReverse,
		},
		{
			name: "empty",
			text: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toSelect(tt.text))
		})
	}
}
