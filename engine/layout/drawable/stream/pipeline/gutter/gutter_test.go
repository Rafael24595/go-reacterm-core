package gutter

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestDrawTransformer(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		size        winsize.Winsize
		lines       []line.Line
		wantLines   []string
		wantFrags   []int
		wantHasNext bool
	}{
		{
			name: "Add left gutter and ignores the right one",
			opts: []Option{
				WithLeftGutter("▌"),
			},
			size: winsize.New(10, 40),
			lines: []line.Line{
				*line.New("golang"),
			},
			wantLines: []string{
				"▌golang",
			},
			wantFrags: []int{
				2,
			},
			wantHasNext: false,
		},
		{
			name: "Add right gutter and ignores the left one",
			opts: []Option{
				WithRightGutter("▐"),
			},
			size: winsize.New(10, 40),
			lines: []line.Line{
				*line.New("golang"),
			},
			wantLines: []string{
				"golang▐",
			},
			wantFrags: []int{
				2,
			},
			wantHasNext: false,
		},
		{
			name: "Add left and right gutter",
			opts: []Option{
				WithGutter(">", "<"),
			},
			size: winsize.New(10, 50),
			lines: []line.Line{
				*line.New("golang"),
				*line.New("ziglang"),
			},
			wantLines: []string{
				">golang<",
				">ziglang<",
			},
			wantFrags: []int{
				3, 3,
			},
			wantHasNext: true,
		},
		{
			name: "Gutter overflow: ignore the gutter and return the line",
			opts: []Option{
				WithLeftGutter("▌ "),
			},
			size: winsize.New(10, 2),
			lines: []line.Line{
				*line.New("golang"),
			},
			wantLines: []string{
				"golang",
			},
			wantFrags: []int{
				1,
			},
			wantHasNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &drawable_test.MockUnit{
				Lines:  tt.lines,
				Status: tt.wantHasNext,
			}

			transform := DrawTransformer(tt.opts...)

			lines, hasNext := transform(tt.size, mock.ToUnit())

			assert.Equal(t, tt.wantHasNext, hasNext)
			for i := range lines {
				assert.Size(t, tt.wantFrags[i], lines[i].Text)
				assert.Equal(t, tt.wantLines[i], text_test.LineToString(lines[i]))
			}
		})
	}
}
