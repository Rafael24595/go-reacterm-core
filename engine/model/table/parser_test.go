package table_test

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
)

type testStruct struct {
	ID     int    `table:"Lang ID"`
	Name   string `table:"Full Name"`
	Date    int
	secret string
}

type emptyStruct struct{}

func TestStructHeaders(t *testing.T) {
	headers := table.StructHeaders[testStruct]()

	assert.Size(t, 3, headers)

	assert.Equal(t, "Lang ID", headers[0])
	assert.Equal(t, "Full Name", headers[1])
	assert.Equal(t, "Date", headers[2])
}

func TestStructFields_WithValues(t *testing.T) {
	s := testStruct{
		ID:     1,
		Name:   "Golang",
		Date:    1257807600,
		secret: "hidden",
	}

	fields := table.StructFields(s)

	assert.Size(t, 3, fields)

	assert.Equal(t, "Lang ID", fields[0].Header)
	assert.Equal(t, 1, fields[0].Value)

	assert.Equal(t, "Full Name", fields[1].Header)
	assert.Equal(t, "Golang", fields[1].Value)

	assert.Equal(t, "Date", fields[2].Header)
	assert.Equal(t, 1257807600, fields[2].Value)
}

func TestStructFields_WithPointer(t *testing.T) {
	s := &testStruct{
		ID:   2,
		Name: "Zig",
		Date:  1454886000,
	}

	fields := table.StructFields(s)

	assert.Size(t, 3, fields)

	assert.Equal(t, "Lang ID", fields[0].Header)
	assert.Equal(t, 2, fields[0].Value)
}

func TestStructFields_WithNilPointer_ShouldReturnEmpty(t *testing.T) {
	var s *testStruct

	fields := table.StructFields(s)

	assert.Empty(t, fields)
}

func TestStructFieds_EmptyStruct_ShouldReturnEmptySlice(t *testing.T) {
	s := emptyStruct{}

	fields := table.StructFields(s)

	assert.Empty(t, fields)
}

func TestStructFieds_NonStruct_ShouldReturnNil(t *testing.T) {
	fields := table.StructFields(123)

	assert.Empty(t, fields)
}

