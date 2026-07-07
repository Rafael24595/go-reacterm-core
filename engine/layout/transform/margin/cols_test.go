package margin

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"
	"github.com/Rafael24595/go-reacterm-core/engine/render/styler"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	render_test "github.com/Rafael24595/go-reacterm-core/test/engine/render"
)

func TestColsLeftTransformer(t *testing.T) {
	styler := styler.NewDefaultSpec()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"  golang"},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 7),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{" golang"},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Cols(
				hint.Fixed(tt.margin),
				cols.WithPosition(style.Right),
			)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Frags(styler, tt.size, result[i].Text))
			}
		})
	}
}

func TestColsRightTransformer(t *testing.T) {
	styler := styler.NewDefaultSpec()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang  "},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 7),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang "},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Cols(
				hint.Fixed(tt.margin),
				cols.WithPosition(style.Left),
			)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Frags(styler, tt.size, result[i].Text))
			}
		})
	}
}

func TestColsCenterTransformer(t *testing.T) {
	styler := styler.NewDefaultSpec()

	tests := []struct {
		name      string
		size      winsize.Winsize
		lines     []string
		margin    winsize.Cols
		wantLen   uint
		wantLines []string
	}{
		{
			name:      "Add all margins when there is enough space.",
			size:      winsize.New(5, 10),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"  golang  "},
		},
		{
			name:      "Add some margins when there is not enough space.",
			size:      winsize.New(5, 8),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{" golang "},
		},
		{
			name:      "Ignore margings when there is not enough space.",
			size:      winsize.New(5, 6),
			lines:     []string{"golang"},
			margin:    2,
			wantLen:   1,
			wantLines: []string{"golang"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transformer := Cols(
				hint.Fixed(tt.margin),
				cols.WithPosition(style.Center),
			)

			lines := make([]text.Line, len(tt.lines))
			for i, l := range tt.lines {
				lines[i] = *text.NewLine(l)
			}

			result := transformer(tt.size, lines)

			assert.Size(t, tt.wantLen, result)

			for i, l := range tt.wantLines {
				assert.Equal(t, l, render_test.Frags(styler, tt.size, result[i].Text))
			}
		})
	}
}
