package margin

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestRowsTopTransformer(t *testing.T) {
	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"golang", "", ""},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(2, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   2,
			wantLines: []string{"golang", ""},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Rows(
				hint.Fixed(tt.margin),
				rows.WithPosition(style.Top),
			)

			lines := make([]line.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = line.FromString(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text_test.LineToString(result[i]))
			}
		})
	}
}

func TestRowsBottomTransformer(t *testing.T) {
	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"", "", "golang"},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(2, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   2,
			wantLines: []string{"", "golang"},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Rows(
				hint.Fixed(tt.margin),
				rows.WithPosition(style.Bottom),
			)

			lines := make([]line.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = line.FromString(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text_test.LineToString(result[i]))
			}
		})
	}
}

func TestRowsMiddleTransformer(t *testing.T) {
	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Rows
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(10, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   5,
			wantLines: []string{"", "", "golang", "", ""},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(3, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   3,
			wantLines: []string{"", "golang", ""},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(1, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Rows(
				hint.Fixed(tt.margin),
				rows.WithPosition(style.Middle),
			)

			lines := make([]line.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = line.FromString(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, text_test.LineToString(result[i]))
			}
		})
	}
}
