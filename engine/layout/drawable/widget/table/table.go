package table

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"

	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/text"
)

const Name = "table_unit"

const defaultMinColSize = winsize.Cols(3 + marker.DefaultElipsisSize)

type TableUnit struct {
	loaded     bool
	lazyLoaded bool
	minColSize winsize.Cols
	table      table.Table
	sections   []section
	cursor     input.MatrixCursor
}

func New(table table.Table, cursor input.MatrixCursor) *TableUnit {
	return &TableUnit{
		loaded:     false,
		lazyLoaded: false,
		minColSize: defaultMinColSize,
		table:      table,
		sections:   make([]section, 0),
		cursor:     cursor,
	}
}

func UnitFromTable(table table.Table, cursor input.MatrixCursor) drawable.Unit {
	return New(table, cursor).ToUnit()
}

func (u *TableUnit) MinCol(size winsize.Cols) *TableUnit {
	u.minColSize = size
	return u
}

func (u *TableUnit) ToUnit() drawable.Unit {
	return drawable.NewBuilder().
		Name(Name).
		Init(u.init).
		Wipe(u.wipe).
		Draw(u.draw).
		ToUnit()
}

func (u *TableUnit) init() {
	u.loaded = true
	u.lazyLoaded = false
}

func (u *TableUnit) wipe() {
	u.lazyLoaded = false
}

func (u *TableUnit) lazyInit(size winsize.Winsize) {
	if u.lazyLoaded {
		return
	}

	u.lazyLoaded = true

	u.sections = newBuilder(u.table, u.cursor).
		setMinSize(u.minColSize).
		render(size)

	for i := range u.sections {
		u.sections[i].header.Drawable.Init()
		u.sections[i].rows.Drawable.Init()
		u.sections[i].footer.Drawable.Init()
	}
}

func (u *TableUnit) draw(size winsize.Winsize) ([]text.Line, bool) {
	assert.True(u.loaded, drawable.MessageInitialized)

	if size.Rows == 0 {
		return make([]text.Line, 0), false
	}

	u.lazyInit(size)

	headers, footers, remaining := u.drawStatic(size)
	bodies, hasNext := u.drawDynamic(size, remaining)

	result := make([]text.Line, size.Rows)
	cursor := 0

	for i, body := range bodies {
		if len(body) == 0 {
			continue
		}

		cursor += copy(result[cursor:], headers[i])
		cursor += copy(result[cursor:], body)
		cursor += copy(result[cursor:], footers[i])
	}

	return result, hasNext
}

func (u *TableUnit) drawStatic(size winsize.Winsize) ([][]text.Line, [][]text.Line, int) {
	headers := make([][]text.Line, len(u.sections))
	footers := make([][]text.Line, len(u.sections))

	remaining := int(size.Rows)
	for i, s := range u.sections {
		header, _ := s.header.Drawable.Draw(size)
		headers[i] = header

		footer, _ := s.footer.Drawable.Draw(size)
		footers[i] = footer

		remaining -= (len(header) + len(footer))
	}

	return headers, footers, remaining
}

func (u *TableUnit) drawDynamic(size winsize.Winsize, remaining int) ([][]text.Line, bool) {
	empty := make(map[int]int)

	sections := len(u.sections)
	if sections == 0 {
		return make([][]text.Line, 0), false
	}

	fixRemaining := remaining - (remaining % sections)

	bodies := make([][]text.Line, sections)
	for fixRemaining > 0 && len(empty) != sections {
		for i, s := range u.sections {
			if fixRemaining <= 0 {
				break
			}

			if _, exists := empty[i]; exists {
				continue
			}

			lines, status := s.rows.Drawable.Draw(size)
			if !status {
				empty[i] = 1
			}

			bodies[i] = append(bodies[i], lines...)

			fixRemaining -= len(lines)
		}
	}

	return bodies, len(empty) != len(u.sections)
}
