package styler

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/format"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
)

func TestJustifyRight_Strict(t *testing.T) {
	spec := spec.JustifyRight(6, "-")

	text := format.TextFromString("hi")

	got := justifyRight(spec, 20, text)

	assert.Equal(t, "----hi", got)
	assert.Equal(t, 6, len(got))
}

func TestJustifyRight_RespectsCols(t *testing.T) {
	spec := spec.JustifyRight(10, "-")

	text := format.TextFromString("hi")

	got := justifyRight(spec, 5, text)

	assert.Equal(t, "---hi", got)
	assert.Equal(t, 5, len(got))
}

func TestJustifyLeft_Strict(t *testing.T) {
	spec := spec.JustifyLeft(6, ".")

	text := format.TextFromString("hi")

	got := justifyLeft(spec, 20, text)

	assert.Equal(t, "hi....", got)
	assert.Equal(t, 6, len(got))
}

func TestJustifyLeft_RespectsCols(t *testing.T) {
	spec := spec.JustifyLeft(10, ".")

	text := format.TextFromString("hi")

	got := justifyLeft(spec, 5, text)

	assert.Equal(t, "hi...", got)
	assert.Equal(t, 5, len(got))
}

func TestJustifyCenter_Strict(t *testing.T) {
	spec := spec.JustifyCenter(6, "-")

	text := format.TextFromString("hi")

	got := justifyCenter(spec, 20, text)

	assert.Equal(t, "--hi--", got)
	assert.Equal(t, 6, len(got))
}

func TestJustifyCenter_RespectsCols(t *testing.T) {
	spec := spec.JustifyCenter(6, "-")

	text := format.TextFromString("hi")

	got := justifyCenter(spec, 4, text)

	assert.Equal(t, "-hi-", got)
	assert.Equal(t, 4, len(got))
}

func TestJustifyCenter_OddSize(t *testing.T) {
	spec := spec.JustifyCenter(7, "-")

	text := format.TextFromString("hi")

	got := justifyCenter(spec, 20, text)

	assert.Equal(t, "--hi---", got)
	assert.Equal(t, 7, len(got))
}

func TestExtendLeft_WithText_Strict(t *testing.T) {
	spec := spec.ExtendLeft(3, "-")

	text := format.TextFromString("hi")

	got := extendLeft(spec, 20, text)

	assert.Equal(t, "---hi", got)
}

func TestExtendLeft_WithoutText_Strict(t *testing.T) {
	spec := spec.ExtendLeft(3)

	text := format.TextFromString("ab")

	got := extendLeft(spec, 20, text)

	assert.Equal(t, "bab", got)
}

func TestExtendRight_WithText_Strict(t *testing.T) {
	spec := spec.ExtendRight(3, "-")

	text := format.TextFromString("hi")

	got := extendRight(spec, 20, text)

	assert.Equal(t, "hi---", got)
}

func TestExtendRight_WithoutText_Strict(t *testing.T) {
	spec := spec.ExtendRight(3)

	text := format.TextFromString("ab")

	got := extendRight(spec, 20, text)

	assert.Equal(t, "aba", got)
}

func TestTruncateLeft_Standard(t *testing.T) {
	tests := []struct {
		name string
		size winsize.Cols
		in   string
		want string
	}{
		{
			name: "keep last 2 characters",
			size: 2,
			in:   "golang",
			want: "ng",
		},
		{
			name: "keep last character",
			size: 1,
			in:   "zig",
			want: "g",
		},
		{
			name: "fallback to minimum 1 when size is 0",
			size: 0,
			in:   "go",
			want: "o",
		},
		{
			name: "handle empty string",
			size: 3,
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := spec.TruncateLeft(tt.size)

			text := format.TextFromString(tt.in)

			got := truncateLeft(spec, text)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && text.Size > 0 {
				assert.Equal(t, tt.size, runes.Measure(got))
			}
		})
	}
}

func TestTruncateLeft_WithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		size     winsize.Cols
		ellipsis string
		in       string
		want     string
	}{
		{
			name:     "prepend ellipsis when space allows",
			size:     5,
			ellipsis: ".",
			in:       "golang",
			want:     "...ng",
		},
		{
			name:     "skip ellipsis if it consumes too much space",
			size:     2,
			ellipsis: ".",
			in:       "ziglang",
			want:     "ng",
		},
		{
			name:     "bypass ellipsis logic when size+elipSize exceeds logical limits",
			size:     1,
			ellipsis: "..",
			in:       "rust",
			want:     "t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := spec.TruncateLeft(tt.size, tt.ellipsis)

			text := format.TextFromString(tt.in)

			got := truncateLeft(spec, text)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && text.Size > 0 {
				assert.Equal(t, tt.size, runes.Measure(got))
			}
		})
	}
}

func TestTruncateRight_Standard(t *testing.T) {
	tests := []struct {
		name string
		size winsize.Cols
		in   string
		want string
	}{
		{
			name: "keep first 2 characters",
			size: 2,
			in:   "golang",
			want: "go",
		},
		{
			name: "keep first character",
			size: 1,
			in:   "ziglang",
			want: "z"},
		{
			name: "fallback to minimum 1 when size is 0",
			size: 0,
			in:   "go",
			want: "g",
		},
		{
			name: "handle empty string",
			size: 2,
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := spec.TruncateRight(tt.size)

			text := format.TextFromString(tt.in)

			got := truncateRight(spec, text)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && text.Size > 0 {
				assert.Equal(t, tt.size, runes.Measure(got))
			}
		})
	}
}

func TestTruncateRight_WithEllipsis(t *testing.T) {
	tests := []struct {
		name     string
		size     winsize.Cols
		ellipsis string
		in       string
		want     string
	}{
		{
			name:     "append ellipsis when space allows",
			size:     5,
			ellipsis: ".",
			in:       "golang",
			want:     "go...",
		},
		{
			name:     "skip ellipsis and return raw trim when space is tight",
			size:     2,
			ellipsis: ".",
			in:       "ziglang",
			want:     "zi",
		},
		{
			name:     "return direct trim when logical size is exceeded",
			size:     1,
			ellipsis: "...",
			in:       "test",
			want:     "t",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := spec.TruncateRight(tt.size, tt.ellipsis)

			text := format.TextFromString(tt.in)

			got := truncateRight(spec, text)

			assert.Equal(t, tt.want, got)

			if tt.size > 0 && text.Size > 0 {
				assert.Equal(t, tt.size, runes.Measure(got))
			}
		})
	}
}

func TestFill_Strict(t *testing.T) {
	text := format.TextFromString("-")
	spec := spec.Fill(10)
	got := fill(spec, 6, text)

	assert.Equal(t, 6, len(got))
	assert.Equal(t, "------", got)
}

func TestFill_Strict_LongText_Even(t *testing.T) {
	text := format.TextFromString("go")

	spec := spec.Fill(20)
	got := fill(spec, 10, text)

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "gogogogogo", got)
}

func TestFill_Strict_LongText_Odd(t *testing.T) {
	text := format.TextFromString("zig")

	spec := spec.Fill(20)
	got := fill(spec, 10, text)

	assert.Equal(t, 10, len(got))
	assert.Equal(t, "zigzigzigz", got)
}
