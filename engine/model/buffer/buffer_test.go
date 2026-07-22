package buffer

import (
	"math"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/model/buffer/processor"
)

func TestRuneBuffer_NumberFilter(t *testing.T) {
	rb := NewRuneBuffer().
		Processor(processor.Number)

	inserted, deleted := rb.Replace([]rune("1A2"), 0, 0)

	assert.Empty(t, deleted)
	
	assert.Equal(t, "12", string(rb.Buffer()))
	assert.Equal(t, "12", string(inserted))
	assert.Equal(t, 2, rb.Size())
}

func TestRuneBuffer_Limit(t *testing.T) {
	handler := processor.Limit(5, processor.Identity)
	rb := NewRuneBuffer().
		Processor(handler)

	rb.Replace([]rune("123"), 0, 0)

	inserted, _ := rb.Replace([]rune("ABCD"), 3, 3)

	assert.Equal(t, "123AB", string(rb.Buffer()))
	assert.Equal(t, "AB", string(inserted))
}

func TestRuneBuffer_SelectionOverwrite(t *testing.T) {
	rb := NewRuneBuffer()
	rb.Replace([]rune("Hello World"), 0, 0)

	inserted, deleted := rb.Replace([]rune("Go"), 6, 11)

	assert.Equal(t, "Hello Go", string(rb.Buffer()))
	assert.Equal(t, "World", string(deleted))
	assert.Equal(t, "Go", string(inserted))
}

func TestRuneBuffer_Slice(t *testing.T) {
	rb := NewRuneBuffer()
	rb.Replace([]rune("0123456789"), 0, 0)

	deleted := rb.Delete(2, 5)

	assert.Equal(t, "0156789", string(rb.Buffer()))
	assert.Equal(t, "234", string(deleted))
}

func TestRuneBuffer_OutOfBounds(t *testing.T) {
	rb := NewRuneBuffer()
	rb.Replace([]rune("Golang"), 0, 0)

	deleted := rb.Delete(2, 100)

	assert.Equal(t, "lang", string(deleted))
}

func TestRuneBuffer_SafeBounds(t *testing.T) {
	rb := NewRuneBuffer()
	rb.Replace([]rune("ABC"), 0, 0)

	_, deleted := rb.Replace([]rune("Z"), 1, 500)

	assert.Equal(t, "BC", string(deleted))
	assert.Equal(t, "AZ", string(rb.Buffer()))
}

func TestRuneBuffer_ReplaceAll(t *testing.T) {
	rb := NewRuneBuffer()
	rb.Replace([]rune("Hello Golang"), 0, 0)

	inserted, deleted := rb.Replace([]rune("Ziglang"), 0, rb.Size())

	assert.Equal(t, "Hello Golang", string(deleted))
	assert.Equal(t, "Ziglang", string(inserted))
	assert.Equal(t, "Ziglang", string(rb.Buffer()))
}

func TestRuneBuffer_VersionStartsAtZero(t *testing.T) {
	buffer := NewRuneBuffer()

	assert.Equal(t, 0, buffer.Version())
}

func TestRuneBuffer_AppendIncrementsVersion(t *testing.T) {
	buffer := NewRuneBuffer()

	before := buffer.Version()

	buffer.Append([]rune("hello"))

	assert.Equal(t, before+1, buffer.Version())
}

func TestRuneBuffer_ReplaceIncrementsVersion(t *testing.T) {
	buffer := NewRuneBuffer()

	buffer.Append([]rune("hello"))

	before := buffer.Version()

	buffer.Replace([]rune("x"), 0, 1)

	assert.Equal(t, before+1, buffer.Version())
}

func TestRuneBuffer_DeleteIncrementsVersion(t *testing.T) {
	buffer := NewRuneBuffer()

	buffer.Append([]rune("hello"))

	before := buffer.Version()

	buffer.Delete(0, 1)

	assert.Equal(t, before+1, buffer.Version())
}

func TestRuneBuffer_CleanIncrementsVersion(t *testing.T) {
	buffer := NewRuneBuffer()

	buffer.Append([]rune("hello"))

	before := buffer.Version()

	buffer.Clean()

	assert.Equal(t, before+1, buffer.Version())
}

func TestRuneBuffer_ReadOperationsDoNotIncrementVersion(t *testing.T) {
	buffer := NewRuneBuffer()

	buffer.Append([]rune("hello"))

	version := buffer.Version()

	_ = buffer.Buffer()
	_ = buffer.Facade()
	_ = buffer.Size()
	_ = buffer.Empty()
	_ = buffer.Range(0, 2)

	assert.Equal(t, version, buffer.Version())
}

func TestRuneBuffer_VersionWrapAround(t *testing.T) {
	buffer := NewRuneBuffer()

	buffer.version = math.MaxUint64

	buffer.Append([]rune("x"))

	assert.Equal(t, 0, buffer.Version())
}
