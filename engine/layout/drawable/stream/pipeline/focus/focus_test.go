package focus

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/action"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestFocusInitTransformer_FocusAtStart(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.LineFromFrags(
				*text.NewFrag("base").
					AddAtom(atom.Focus),
				*text.NewFrag("_"),
				*text.NewFrag("02"),
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
	assert.Equal(t, "base_01", text.LineToString(&lines[0]))
	assert.Equal(t, "base_02", text.LineToString(&lines[1]))
}

func TestFocusInitTransformer_FocusAtEnd(t *testing.T) {
	mock := &drawable_test.MockUnit{
		Lines: []text.Line{
			*text.NewLine("base_01"),
			*text.NewLine("base_02"),
			*text.LineFromFrags(
				*text.NewFrag("base"),
				*text.NewFrag("_").
					AddAtom(atom.Focus),
				*text.NewFrag("03"),
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
	assert.Equal(t, "base_02", text.LineToString(&lines[0]))
	assert.Equal(t, "base_03", text.LineToString(&lines[1]))
}
