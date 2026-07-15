package table

import (
	"github.com/Rafael24595/go-reacterm-core/engine/commons/structure/heap"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/runes"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/primitive/lines"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/isolated"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/atom"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/frag"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

type cellRenderer func(col int, header string) frag.Frag

type section struct {
	header drawable.Unit
	rows   drawable.Unit
	footer drawable.Unit
}
type col struct {
	name string
	size winsize.Cols
}

type builder struct {
	minSize winsize.Cols
	table   table.Table
	cursor  input.MatrixCursor
}

func newBuilder(table table.Table, cursor input.MatrixCursor) *builder {
	return &builder{
		minSize: winsize.Cols(defaultMinColSize),
		table:   table,
		cursor:  cursor,
	}
}

func (b *builder) setMinSize(size winsize.Cols) *builder {
	b.minSize = size
	return b
}

func (b builder) render(size winsize.Winsize) []section {
	sections := make([]section, 0)

	separator := b.table.GetSeparator()
	headers := b.table.GetHeaders()

	chunks := b.maxColsChunks(size)

	for _, chunk := range chunks {
		headers, fixedCursor := b.filterHeaders(headers, chunk)

		specCover := spec.ExtendRight(
			b.calcRowCapacity(chunk),
		)

		topCover := line.FromFrags(
			*frag.New(separator.Top).AddSpec(specCover),
		)

		bottomCover := line.FromFrags(
			*frag.New(separator.Bottom).AddSpec(specCover),
		)

		bodyRows := b.renderBody(chunk, headers, separator, fixedCursor)
		if len(bodyRows) == 0 {
			continue
		}

		headerRow := b.renderHeaders(chunk, headers, separator)

		sections = append(sections, section{
			header: isolated.UnitFromLines(
				*topCover, *headerRow, *topCover,
			),
			rows: lines.UnitFromLines(bodyRows...),
			footer: isolated.UnitFromLines(
				*bottomCover,
			),
		})
	}

	return sections
}

func (b builder) maxColsChunks(size winsize.Winsize) []table.MaxCols {
	maxCols, fits := b.evalMaxCols(size)
	if fits {
		return []table.MaxCols{maxCols}
	}

	return b.splitMaxCols(maxCols, size)
}

func (b builder) evalMaxCols(size winsize.Winsize) (table.MaxCols, bool) {
	headers := b.table.GetHeaders()
	maxCols := b.table.MaxCols()

	capacity := b.calcRowCapacity(maxCols)
	if capacity <= size.Cols {
		return maxCols, true
	}

	heap := heap.NewMaxBy(func(c col) winsize.Cols {
		return c.size
	})

	for _, header := range headers {
		heap.Push(col{header, maxCols[header]})
	}

	excess := capacity.Sub(size.Cols)

	for excess > 0 {
		col, ok := heap.Peek()
		if !ok || col.size <= b.minSize {
			break
		}

		col, _ = heap.Pop()

		col.size = col.size.Sub(1)
		excess = excess.Sub(1)

		heap.Push(col)
	}

	fixedMaxCols := make(table.MaxCols)
	for heap.Len() > 0 {
		c, _ := heap.Pop()
		fixedMaxCols[c.name] = c.size
	}

	return fixedMaxCols, excess == 0
}

func (b builder) splitMaxCols(maxCols table.MaxCols, size winsize.Winsize) []table.MaxCols {
	chunks := make([]table.MaxCols, 0)

	separator := b.table.GetSeparator()
	headers := b.table.GetHeaders()

	leftLen := runes.Measure(separator.Left)
	centerLen := runes.Measure(separator.Center)
	rightLen := runes.Measure(separator.Right)

	headersLen := len(headers)

	chunk := make(table.MaxCols)
	count := leftLen

	flush := func() {
		if len(chunk) == 0 {
			return
		}

		chunks = append(chunks, chunk)

		chunk = make(table.MaxCols)
		count = leftLen
	}

	for i, header := range headers {
		measure := min(maxCols[header], size.Cols)

		needed := count + measure + rightLen
		if needed > size.Cols {
			flush()
		}

		chunk[header] = measure
		count += measure

		if i < headersLen-1 {
			count += centerLen
		}
	}

	flush()

	return chunks
}

func (b builder) calcRowCapacity(maxCols table.MaxCols) winsize.Cols {
	cols := max(0, len(maxCols)-1)

	separator := b.table.GetSeparator()

	centerMeasure := runes.Measure(separator.Center)
	leftMeasure := runes.Measure(separator.Left)
	rightMeasure := runes.Measure(separator.Right)

	joinMeasure := winsize.Cols(cols) * centerMeasure
	borderMeasure := leftMeasure + rightMeasure

	total := winsize.Cols(0)
	for _, c := range maxCols {
		total += c
	}

	return total + joinMeasure + borderMeasure
}

func (b builder) filterHeaders(
	headers []string,
	maxCols table.MaxCols,
) ([]string, *input.MatrixCursor) {
	cursor := int(b.cursor.Col)

	var fixCursor *input.MatrixCursor
	var fixedCol uint16

	filtered := make([]string, 0)
	for i, header := range headers {
		if _, ok := maxCols[header]; !ok {
			continue
		}

		filtered = append(filtered, header)
		if i == cursor {
			fixCursor = input.NewMatrixCursor(
				b.cursor.Row,
				fixedCol,
				b.cursor.Show,
			)
		}

		fixedCol += 1
	}

	return filtered, fixCursor
}

func (b builder) renderHeaders(
	maxCols table.MaxCols,
	headers []string,
	separator marker.TableSeparatorMeta,
) *line.Line {
	renderer := func(col int, header string) frag.Frag {
		maxCol := maxCols[header]

		spec := spec.Merge(
			spec.JustifyCenter(maxCol),
			spec.TruncateRight(maxCol, marker.DefaultElipsisText),
		)

		return *frag.New(header).AddSpec(spec)
	}

	return b.renderRow(headers, separator, renderer)
}

func (b builder) renderBody(
	maxCols table.MaxCols,
	headers []string,
	separator marker.TableSeparatorMeta,
	cursor *input.MatrixCursor,
) []line.Line {
	rows := table.RowCount(
		headers, b.table.GetColumns(),
	)

	lines := make([]line.Line, rows)

	for row := range rows {
		renderer := func(col int, header string) frag.Frag {
			maxCol := maxCols[header]
			return *b.renderCell(
				maxCol, cursor, header, row, uint16(col),
			)
		}

		lines[row] = *b.renderRow(headers, separator, renderer)
	}

	return lines
}

func (b builder) renderRow(
	headers []string,
	separator marker.TableSeparatorMeta,
	cellRenderer cellRenderer,
) *line.Line {
	headersLen := len(headers)
	capacity := 2*headersLen + 1

	frags := make([]frag.Frag, 0, capacity)

	frags = append(frags,
		*frag.New(separator.Left).
			AddAtom(atom.Wrap),
	)

	for col, header := range headers {
		frags = append(frags,
			cellRenderer(col, header),
		)

		if col < headersLen-1 {
			frags = append(frags,
				*frag.New(separator.Center).
					AddAtom(atom.Wrap),
			)
		}
	}

	frags = append(frags,
		*frag.New(separator.Right).
			AddAtom(atom.Wrap),
	)

	return line.FromFrags(frags...)
}

func (b builder) renderCell(
	maxCol winsize.Cols,
	cursor *input.MatrixCursor,
	header string,
	row uint16,
	col uint16,
) *frag.Frag {
	atm := atom.Wrap
	if cursor != nil && cursor.IsAt(row, col) {
		atm = atom.Merge(atm, atom.Select, atom.Focus)
	}

	cell, ok := b.table.FindCell(header, row)
	if !ok {
		scp := spec.ExtendRight(maxCol)

		return frag.New("").
			AddSpec(scp).
			AddAtom(atm)
	}

	scp := spec.Merge(
		spec.JustifyLeft(maxCol),
		spec.TruncateRight(maxCol, marker.DefaultElipsisText),
	)

	return frag.New(cell).
		AddSpec(scp).
		AddAtom(atm)
}
