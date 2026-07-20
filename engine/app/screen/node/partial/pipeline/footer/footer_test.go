package footer

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/node/partial/pipeline"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/spatial/stack"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"

	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

func TestFooter_InsertsBefore(t *testing.T) {
	vm := viewmodel.ViewModel{
		Footer: stack.NewVStack(
			lines.UnitFromLines(
				*line.New("line_01"),
			),
		),
	}

	units := vm.Footer.Units()
	assert.Size(t, 1, units)

	line := line.New("line_02")
	transformer := Transformer(pipeline.Before, *line)
	vm = transformer(vm)

	units = vm.Footer.Units()
	assert.Size(t, 2, units)

	unit := units[0]

	unit.Drawable.Boot()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 1,
		Cols: 10,
	})

	assert.Equal(t, "line_02", text_test.LineToString(lines[0]))
}

func TestFooter_InsertsAfter(t *testing.T) {
	vm := viewmodel.ViewModel{
		Footer: stack.NewVStack(
			lines.UnitFromLines(
				*line.New("line_01"),
			),
		),
	}

	units := vm.Footer.Units()
	assert.Size(t, 1, units)

	line := line.New("line_02")
	transformer := Transformer(pipeline.After, *line)
	vm = transformer(vm)

	units = vm.Footer.Units()
	assert.Size(t, 2, units)

	unit := units[1]

	unit.Drawable.Boot()
	lines, _ := unit.Drawable.Draw(winsize.Winsize{
		Rows: 1,
		Cols: 10,
	})

	assert.Equal(t, "line_02", text_test.LineToString(lines[0]))
}
