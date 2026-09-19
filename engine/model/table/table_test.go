package table

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

func TestNewTable_ShouldInitializeEmptyTable(t *testing.T) {
	tbl := NewTable()

	assert.Equal(t, 0, tbl.ColCount())
	assert.Equal(t, 0, tbl.RowCount())
	assert.Equal(t, marker.DefaultTableSeparator, tbl.Separator())
}

func TestAddHeaders_ShouldAddHeadersWithoutDuplicates(t *testing.T) {
	tbl := NewTable()

	tbl.AddHeaders("ID", "Lang")
	tbl.AddHeaders("Lang", "Age")

	headers := tbl.Headers()

	assert.Size(t, 3, headers)

	assert.Equal(t, "ID", headers[0])
	assert.Equal(t, "Lang", headers[1])
	assert.Equal(t, "Age", headers[2])
}

func TestField_ShouldExpandRowsDynamically(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")

	tbl.SetCell("Name", 2, "Golang")

	col := tbl.Columns()["Name"]

	assert.Size(t, 3, col)
	assert.Equal(t, "", col[0])
	assert.Equal(t, "", col[1])
	assert.Equal(t, "Golang", col[2])
}

func TestField_WithInvalidHeader_ShouldDoNothing(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("ID")

	tbl.SetCell("Invalid", 0, "X")

	assert.Empty(t, tbl.Columns()["ID"])
}

func TestSize_ShouldCalculateMaxWidth(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")

	tbl.SetCell("Name", 0, "zig")
	tbl.SetCell("Name", 1, "golang")

	maxCols := tbl.MeasureColWidths()

	assert.Size(t, maxCols["Name"], "golang")
}

func TestSize_ShouldConsiderHeaderLength(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("VeryLongHeader")

	tbl.SetCell("VeryLongHeader", 0, "go")

	maxCols := tbl.MeasureColWidths()

	assert.Size(t, maxCols["VeryLongHeader"], "VeryLongHeader")
}

func TestCols_ShouldReturnHeaderCount(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("A", "B", "C")

	assert.Equal(t, 3, tbl.ColCount())
}

func TestRows_ShouldReturnMaxRowCount(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("A", "B")

	tbl.SetCell("A", 0, "x")
	tbl.SetCell("B", 2, "y")

	assert.Equal(t, 3, tbl.RowCount())
}

func TestWithSeparator_ShouldOverrideDefault(t *testing.T) {
	tbl := NewTable()

	sep := marker.TableSeparatorMeta{
		Top:    "=",
		Bottom: "=",
		Center: "::",
		Left:   "[",
		Right:  "]",
	}

	ret := tbl.WithSeparator(sep)

	assert.Equal(t, sep, tbl.Separator())
	assert.Equal(t, ret, tbl)
}

func TestFindCell_ShouldReturnContentWhenCellExists(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")
	tbl.SetCell("Name", 0, "Golang")

	val, ok := tbl.FindCell("Name", 0)

	assert.True(t, ok)
	assert.Equal(t, "Golang", val)
}

func TestFindCell_WithInvalidHeader_ShouldReturnFalse(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")
	tbl.SetCell("Name", 0, "Golang")

	val, ok := tbl.FindCell("Invalid", 0)

	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestFindCell_WithRowOutOfBounds_ShouldReturnFalse(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")
	tbl.SetCell("Name", 0, "Golang")

	val, ok := tbl.FindCell("Name", 1)

	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestFindCell_WithDynamicallyExpandedCell_ShouldReturnEmptyStringAndTrue(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")

	tbl.SetCell("Name", 2, "Zig")

	val, ok := tbl.FindCell("Name", 1)

	assert.True(t, ok)
	assert.Equal(t, "", val)
}

func TestFindCellByCoords_ShouldReturnContent(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("ID", "Name")
	tbl.SetCell("Name", 0, "Go")

	val, ok := tbl.FindCellByCoords(0, 1)

	assert.True(t, ok)
	assert.Equal(t, "Go", val)
}

func TestFindCellByCoords_OutOfBounds_ShouldReturnFalse(t *testing.T) {
	tbl := NewTable()
	tbl.AddHeaders("Name")
	tbl.SetCell("Name", 0, "Go")

	_, okCol := tbl.FindCellByCoords(0, 5)
	assert.False(t, okCol)

	_, okRow := tbl.FindCellByCoords(1, 0)
	assert.False(t, okRow)
}

func TestRowCount_Empty(t *testing.T) {
	headers := []string{"ID", "Name"}
	cols := map[string][]string{
		"ID":   {},
		"Name": {},
	}

	assert.Equal(t, 0, RowCount(headers, cols))
}

func TestRowCount_UnevenColumns(t *testing.T) {
	headers := []string{"ID", "Name", "Date"}
	cols := map[string][]string{
		"ID":   {"1", "2"},
		"Name": {"Golang", "Zig", "Rust"},
		"Date": {"1257807600"},
	}

	assert.Equal(t, 3, RowCount(headers, cols))
}
