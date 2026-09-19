package table

import (
	"fmt"
	"slices"

	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
)

// ColWidths maps header names to their maximum calculated visual column width.
type ColWidths map[string]winsize.Cols

// Table represents a dynamic grid of string data structured by column headers.
type Table struct {
	cols      map[string][]string
	headers   []string
	separator marker.TableSeparatorMeta
}

// NewTable initializes and returns an empty Table with default separator styling.
func NewTable() *Table {
	return &Table{
		headers:   make([]string, 0),
		cols:      make(map[string][]string),
		separator: marker.DefaultTableSeparator,
	}
}

// Separator returns the current table separator metadata.
func (t *Table) Separator() marker.TableSeparatorMeta {
	return t.separator
}

// WithSeparator configures custom separator metadata for rendering.
func (t *Table) WithSeparator(separator marker.TableSeparatorMeta) *Table {
	t.separator = separator
	return t
}

// Headers returns the ordered list of table column headers.
func (t *Table) Headers() []string {
	return t.headers
}

// AddHeaders appends unique headers to the table, preserving insertion order.
func (t *Table) AddHeaders(headers ...string) *Table {
	for _, v := range headers {
		if slices.Contains(t.headers, v) {
			continue
		}

		t.headers = append(t.headers, v)
		t.cols[v] = make([]string, 0)
	}

	return t
}

// Columns returns the internal map of column headers to cell contents.
func (t *Table) Columns() map[string][]string {
	return t.cols
}

// FindCell retrieves cell content by header name and row index.
// Returns false if the header does not exist or the row is out of bounds.
func (t *Table) FindCell(header string, row uint16) (string, bool) {
	col, ok := t.cols[header]
	if !ok || row >= uint16(len(col)) {
		return "", false
	}

	return col[row], true
}

// FindCellByCoords retrieves cell content by zero-based column and row coordinates.
// Returns false if coordinates are out of bounds.
func (t *Table) FindCellByCoords(row, col uint16) (string, bool) {
	if col >= uint16(len(t.headers)) {
		return "", false
	}
	return t.FindCell(t.headers[col], row)
}

// SetCell assigns data formatted as a string to a specific header and row, expanding column height if necessary.
func (t *Table) SetCell(header string, row uint16, data any) *Table {
	col, ok := t.cols[header]
	if !ok {
		return t
	}

	colLen := uint16(len(col))
	if row >= colLen {
		for i := colLen; i <= row; i++ {
			col = append(col, "")
		}
	}

	col[row] = fmt.Sprintf("%v", data)
	t.cols[header] = col

	return t
}

// MeasureColWidths computes the maximum visual rune width required for each column (header and cells).
func (t *Table) MeasureColWidths() ColWidths {
	size := make(ColWidths)
	for _, h := range t.headers {
		if _, ok := size[h]; !ok {
			size[h] = runes.MeasureCols(h)
		}

		for _, c := range t.cols[h] {
			size[h] = max(size[h], runes.MeasureCols(c))
		}
	}

	return size
}

// ColCount returns the total number of table headers.
func (t *Table) ColCount() uint16 {
	return uint16(len(t.headers))
}

// RowCount returns the maximum row height across all columns.
func (t *Table) RowCount() uint16 {
	return RowCount(t.headers, t.cols)
}

// RowCount calculates the maximum row height across all specified column headers.
func RowCount(headers []string, cols map[string][]string) uint16 {
	maxRows := 0
	for _, h := range headers {
		maxRows = max(
			maxRows, len(cols[h]),
		)
	}
	return uint16(maxRows)
}
