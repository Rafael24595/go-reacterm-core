package table

import (
	"strings"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/marker"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style/spec"
	
	text_test "github.com/Rafael24595/go-reacterm-core/test/engine/render/text"
)

var separator = marker.TableSeparatorMeta{
	Center: "|",
	Left:   "|",
	Right:  "|",
}

func TestBuilder_FilterHeaders_FiltersCorrectly(t *testing.T) {
	maxCols := table.MaxCols{
		"id":   5,
		"name": 10,
	}

	headers := []string{"id", "name", "email"}

	result, _ := builder{}.filterHeaders(headers, maxCols)

	assert.Size(t, 2, result)
	assert.Equal(t, "id", result[0])
	assert.Equal(t, "name", result[1])
}

func TestBuilder_RenderHeaders_Basic(t *testing.T) {
	maxCols := table.MaxCols{
		"id":   4,
		"name": 6,
	}

	headers := []string{"id", "name"}

	line := builder{}.renderHeaders(maxCols, headers, separator)

	assert.Equal(t, "|id|name|", text_test.LineToString(*line))
}

func TestBuilder_RenderHeaders_Structure(t *testing.T) {
	maxCols := table.MaxCols{
		"id":   4,
		"name": 6,
	}

	headers := []string{"id", "name"}

	sep := marker.TableSeparatorMeta{
		Left:   "|",
		Center: "|",
		Right:  "|",
	}

	line := builder{}.renderHeaders(maxCols, headers, sep)

	wantFrags := 2*len(headers) + 1

	assert.Size(t, wantFrags, line.Text)

	assert.Equal(t, "|", line.Text[0].Text())
	assert.Equal(t, "id", line.Text[1].Text())
	assert.Equal(t, "|", line.Text[2].Text())
	assert.Equal(t, "name", line.Text[3].Text())
	assert.Equal(t, "|", line.Text[4].Text())

	assert.NotEqual(t, spec.KindNone, line.Text[1].Spec().Kind())
	assert.NotEqual(t, spec.KindNone, line.Text[3].Spec().Kind())
}

func TestBuilder_RenderBody_Basic(t *testing.T) {
	maxCols := table.MaxCols{
		"id":   4,
		"name": 6,
	}

	headers := []string{"id", "name"}

	table := table.NewTable().
		SetHeaders(headers...).
		SetCell("id", 0, "1").
		SetCell("id", 1, "2").
		SetCell("name", 0, "golang").
		SetCell("name", 1, "ziglang")

	lines := builder{table: *table}.renderBody(
		maxCols, headers, separator, &input.MatrixCursor{},
	)

	assert.Size(t, 2, lines)

	assert.Equal(t, "|1|golang|", text_test.LineToString(lines[0]))
	assert.Equal(t, "|2|ziglang|", text_test.LineToString(lines[1]))
}

func TestBuilder_EvalMaxCols_NoReductionNeeded(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B").
		SetCell("A", 0, strings.Repeat("0", 5)).
		SetCell("B", 0, strings.Repeat("0", 5))

	builder := builder{
		table: *table,
	}

	result, status := builder.evalMaxCols(winsize.Winsize{
		Cols: 20,
	})

	assert.True(t, status)

	assert.Equal(t, 5, result["A"])
	assert.Equal(t, 5, result["B"])
}

func TestBuilder_EvalMaxCols_ReducesLargestColumn(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B").
		SetCell("A", 0, strings.Repeat("0", 10)).
		SetCell("B", 0, strings.Repeat("0", 5))

	builder := builder{
		table: *table,
	}

	result, status := builder.evalMaxCols(winsize.Winsize{
		Cols: 14,
	})

	assert.True(t, status)

	assert.GreaterOrEqual(t, 14, builder.calcRowCapacity(result))
	assert.LessThan(t, 10, result["A"])
}

func TestBuilder_EvalMaxCols_RespectsMinWidth(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B").
		SetCell("A", 0, strings.Repeat("0", 4)).
		SetCell("B", 0, strings.Repeat("0", 4))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result, status := builder.evalMaxCols(winsize.Winsize{
		Cols: 5,
	})

	assert.False(t, status)

	assert.GreaterOrEqual(t, 4, result["A"])
	assert.GreaterOrEqual(t, 4, result["B"])
}

func TestBuilder_EvalMaxCols_ExactFit(t *testing.T) {
	aMeasure := winsize.Cols(8)
	bMeasure := winsize.Cols(7)

	sepMeasure := winsize.Cols(3)

	totalMeasure := aMeasure + bMeasure + sepMeasure

	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B").
		SetCell("A", 0, strings.Repeat("0", int(aMeasure))).
		SetCell("B", 0, strings.Repeat("0", int(bMeasure)))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result, status := builder.evalMaxCols(winsize.Winsize{
		Cols: totalMeasure,
	})

	assert.True(t, status)

	assert.Equal(t, totalMeasure, builder.calcRowCapacity(result))
}

func TestBuilder_EvalMaxCols_MultipleColumnsReduction(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B", "C").
		SetCell("A", 0, strings.Repeat("0", 10)).
		SetCell("B", 0, strings.Repeat("0", 9)).
		SetCell("C", 0, strings.Repeat("0", 8))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result, status := builder.evalMaxCols(winsize.Winsize{
		Cols: 22,
	})

	assert.True(t, status)

	assert.Equal(t, 22, builder.calcRowCapacity(result))

	assert.NotEqual(t, 10, result["A"])
	assert.NotEqual(t, 9, result["B"])
	assert.NotEqual(t, 8, result["C"])
}

func TestBuilder_MaxColsChunks_FitsInOne(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B", "C").
		SetCell("A", 0, strings.Repeat("0", 10)).
		SetCell("B", 0, strings.Repeat("0", 20)).
		SetCell("C", 0, strings.Repeat("0", 10))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result := builder.maxColsChunks(winsize.Winsize{
		Cols: 50,
	})

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 10, result[0]["A"])
	assert.Equal(t, 20, result[0]["B"])
	assert.Equal(t, 10, result[0]["C"])
}

func TestBuilder_MaxColsChunks_MustSplit(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("A", "B", "C", "D").
		SetCell("A", 0, strings.Repeat("0", 20)).
		SetCell("B", 0, strings.Repeat("0", 10)).
		SetCell("C", 0, strings.Repeat("0", 15)).
		SetCell("D", 0, strings.Repeat("0", 15))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result := builder.maxColsChunks(winsize.Winsize{
		Cols: 25,
	})

	assert.True(t, len(result) > 1)

	for _, table := range result {
		total := winsize.Cols(0)
		for _, v := range table {
			total += v
		}
		assert.True(t, total <= 25)
	}
}

func TestBuilder_MaxColsChunks_ColumnWiderThanTerminal(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders("XL").
		SetCell("XL", 0, strings.Repeat("0", 100))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result := builder.maxColsChunks(winsize.Winsize{
		Cols: 80,
	})

	leftSepMeasure := winsize.Cols(1)
	rightSepMeasure := winsize.Cols(1)

	sepMeasure := leftSepMeasure + rightSepMeasure

	assert.Equal(t, 1, len(result))
	assert.Equal(t, 80-sepMeasure, result[0]["XL"])
}

func TestBuilder_MaxColsChunks_EmptyMap(t *testing.T) {
	table := table.NewTable().
		SetSeparator(separator)

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	result := builder.maxColsChunks(winsize.Winsize{})

	assert.Equal(t, 0, len(result))
}

func TestAdjustSize_Deterministic(t *testing.T) {
	headers := []string{"A", "B", "C", "D", "E", "F", "G"}

	table := table.NewTable().
		SetSeparator(separator).
		SetHeaders(headers...).
		SetCell("A", 0, strings.Repeat("0", 12)).
		SetCell("B", 0, strings.Repeat("0", 10)).
		SetCell("C", 0, strings.Repeat("0", 8)).
		SetCell("D", 0, strings.Repeat("0", 6)).
		SetCell("E", 0, strings.Repeat("0", 10)).
		SetCell("F", 0, strings.Repeat("0", 8)).
		SetCell("G", 0, strings.Repeat("0", 12))

	builder := builder{
		table:   *table,
		minSize: defaultMinColSize,
	}

	var firstResult map[string]winsize.Cols

	for i := range 100 {
		result, status := builder.evalMaxCols(winsize.Winsize{
			Cols: 50,
		})

		assert.True(t, status)

		if i == 0 {
			firstResult = result
			continue
		}

		for _, header := range headers {
			assert.Equal(t, firstResult[header], result[header])
		}
	}
}
