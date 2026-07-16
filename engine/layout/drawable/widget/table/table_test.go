package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"

	drawable_test "github.com/Rafael24595/go-reacterm-core/test/engine/layout/drawable"
)

func TestTable_UnitBasicSuite(t *testing.T) {
	unit := UnitFromTable(
		*table.NewTable(),
		*input.NewMatrixCursor(0, 0, false),
	)
	drawable_test.Test_UnitBasicSuite(t, unit)
}

func TestTable_LazyBoot(t *testing.T) {
	unit := New(
		*table.NewTable().
			SetHeaders("lang").
			SetCell("lang", 0, "golang"),
		*input.NewMatrixCursor(0, 0, false),
	)

	assert.Empty(t, unit.sections)

	unit.boot()
	unit.draw(winsize.Winsize{
		Rows: 3,
		Cols: 11,
	})

	assert.Size(t, 1, unit.sections)
}
