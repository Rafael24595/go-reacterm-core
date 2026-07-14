package focus

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestFocusInitTransformer_FocusAtStart(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.LineFromFrags(
				*frag.New("base").
					AddAtom(atom.Focus),
				*frag.New("_"),
				*frag.New("02"),
			),
			*text.NewLine("base_03"),
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
	assert.Equal(t, "base_01", text_test.LineToString(&lines[0]))
	assert.Equal(t, "base_02", text_test.LineToString(&lines[1]))
}

func TestFocusInitTransformer_FocusAtEnd(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
			*text.LineFromFrags(
				*frag.New("base"),
				*frag.New("_").
					AddAtom(atom.Focus),
				*frag.New("03"),
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
	assert.Equal(t, "base_02", text_test.LineToString(&lines[0]))
	assert.Equal(t, "base_03", text_test.LineToString(&lines[1]))
}
