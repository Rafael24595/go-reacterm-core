package composer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/layer"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestStandard_FixedAndPaged(t *testing.T) {
	size := winsize.Winsize{Rows: 6, Cols: 10}

	vm := viewmodel.New()

	vm.Header.Push(
		drain.UnitFromLines(
			line.TextSpec("HEADER", spec.AlignRight()),
		),
	)

	vm.Kernel.Push(
		lines.UnitFromLines(
			line.TextSpec("=", spec.Cover()),
			line.TextSpec("LINE TWO", spec.AlignRight()),
			line.TextSpec("LINE THREE IS LONG", spec.AlignRight()),
			line.TextSpec("LINE FOUR", spec.AlignRight()),
		),
	)

	frg := frag.FromStrings("INPUT")
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromFrags(frg...),
		},
	}

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	state := &state.UIState{}

	_, lines := Standard(state, *vm, size)

	assert.Size(t, int(size.Rows), lines)
	assert.Equal(t, "HEADER", lines[0].Text[0].Text())

	inputLine := lines[len(lines)-1]
	expectedInput := "> INPUT"

	assert.Equal(t, expectedInput, text_test.FragsToString(inputLine.Text))

	for i := 1; i < len(lines)-1; i++ {
		width := winsize.Cols(0)
		for _, f := range lines[i].Text {
			width += frag.Measure(size.Cols, f)
		}

		assert.LessOrEqual(t, size.Cols, width)
	}
}

func TestStandard_InitializeLayers(t *testing.T) {
	size := winsize.Winsize{Rows: 5, Cols: 8}

	uiState := state.NewUIState()

	vm := viewmodel.New()

	vm.Header.PushWithOpts(
		drain.UnitFromLines(
			line.TextSpec("golang", spec.AlignRight()),
		),
		layer.Fixed[winsize.Rows](1),
	)
	vm.Kernel.PushWithOpts(
		lines.UnitFromLines(
			line.TextSpec("rust", spec.AlignRight()),
		),
		layer.Fixed[winsize.Rows](1),
	)
	vm.Footer.PushWithOpts(
		drain.UnitFromLines(
			line.TextSpec("Ziglang", spec.AlignRight()),
		),
		layer.Fixed[winsize.Rows](1),
	)

	frag := frag.FromStrings("X")
	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromFrags(frag...),
		},
	}

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Kernel.HasNext())
	assert.True(t, vm.Footer.HasNext())

	Standard(uiState, *vm, size)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Kernel.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
