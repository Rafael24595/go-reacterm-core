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
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestStandard_FixedAndPaged(t *testing.T) {
	ws := winsize.New(6, 10)
	us := state.NewUIState()
	vm := viewmodel.New()

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromString("INPUT"),
		},
	}

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

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	_, lines := Standard(us, ws, *vm)

	assert.Size(t, ws.Rows, lines)
	assert.Equal(t, "HEADER", lines[0].AtOrZero(0).Text())

	inputLine := lines[len(lines)-1]
	expectedInput := "> INPUT"

	assert.Equal(t, expectedInput, text_test.LineToString(inputLine))

	for i := range len(lines) {
		width := line.Measure(lines[i], 0)
		assert.LessOrEqual(t, ws.Cols, width)
	}
}

func TestStandard_InitializeLayers(t *testing.T) {
	ws := winsize.New(5, 8)
	us := state.NewUIState()
	vm := viewmodel.New()

	mock := &drawable_test.MockUnit{
		Status: false,
		Lines: []line.Line{
			line.FromString("X"),
		},
	}

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

	vm.Footer.Unshift(
		inputline.Wrap(
			mock.ToUnit(),
		),
	)

	assert.True(t, vm.Header.HasNext())
	assert.True(t, vm.Kernel.HasNext())
	assert.True(t, vm.Footer.HasNext())

	Standard(us, ws, *vm)

	assert.False(t, vm.Header.HasNext())
	assert.False(t, vm.Kernel.HasNext())
	assert.False(t, vm.Footer.HasNext())
}
