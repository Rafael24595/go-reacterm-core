package focus

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestFocusInitTransformer_FocusAtStart(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("base_01"),
			*line.FromFrags(
				frag.TextAtom("base", atom.Focus),
				frag.FromString("_"),
				frag.FromString("02"),
			),
			*line.New("base_03"),
		},
		Status: true,
	}

	transformer := DrawTransformer(
		action.Paged(),
	)

	lines, status := transformer(winsize.Winsize{
		Rows: 2,
		Cols: 10,
	}, mock.ToUnit())

	assert.Size(t, 2, lines)

	assert.False(t, status)
	assert.Equal(t, "base_01", text_test.LineToString(lines[0]))
	assert.Equal(t, "base_02", text_test.LineToString(lines[1]))
}

func TestFocusInitTransformer_FocusAtEnd(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []line.Line{
			*line.New("base_01"),
			*line.New("base_02"),
			*line.FromFrags(
				frag.FromString("base"),
				frag.TextAtom("_", atom.Focus),
				frag.FromString("03"),
			),
		},
		Batch: 1,
	}

	transformer := DrawTransformer(
		action.Scroll(),
	)

	lines, status := transformer(winsize.Winsize{
		Rows: 2,
		Cols: 10,
	}, mock.ToUnit())

	assert.Size(t, 2, lines)

	assert.False(t, status)
	assert.Equal(t, "base_02", text_test.LineToString(lines[0]))
	assert.Equal(t, "base_03", text_test.LineToString(lines[1]))
}
