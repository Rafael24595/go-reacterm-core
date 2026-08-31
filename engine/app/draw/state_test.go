package draw

import (
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/render/text/line"
)

func TestNewState_DefaultInitialization(t *testing.T) {
	st := NewState()

	assert.Empty(t, st.Buffer)
	assert.Equal(t, 0, st.Cursor)
	assert.Equal(t, 0, st.Page)
	assert.False(t, st.Focus)
	assert.NotNil(t, st.Work)
}

func TestNewState_WithCustomSize(t *testing.T) {
	st := NewState(10)

	assert.Size(t, 10, st.Buffer)
	assert.Equal(t, 0, st.Cursor)
}

func TestState_WriteLine_AppendsAndAdvancesCursor(t *testing.T) {
	st := NewState(2)

	st.WriteLine(line.Line{})

	assert.Equal(t, 1, st.Cursor)

	st.WriteLine(line.Line{})

	assert.Equal(t, 2, st.Cursor)
	assert.True(t, st.IsFull())
}

func TestState_Written_ReturnsOnlyWrittenLines(t *testing.T) {
	st := NewState(5).
		WriteLine(line.Line{}).
		WriteLine(line.Line{})

	written := st.Written()
	assert.Size(t, 2, written)
}

func TestState_EnableFocus_AccumulatesFocusState(t *testing.T) {
	st := NewState()

	st.EnableFocus(false)
	assert.False(t, st.Focus)

	st.EnableFocus(true)
	assert.True(t, st.Focus)

	st.EnableFocus(false)
	assert.True(t, st.Focus)
}

func TestState_IsPaginated(t *testing.T) {
	st := NewState()
	assert.False(t, st.IsPaginated())

	st.Page = 1
	assert.True(t, st.IsPaginated())
}

func TestState_Reset_ClearsBufferAndResetsState(t *testing.T) {
	st := NewState(3).
		EnableFocus(true).
		WriteLine(line.Line{}).
		WriteLine(line.Line{})

	st.Reset()

	assert.Equal(t, 0, st.Cursor)
	assert.False(t, st.Focus)
	assert.Empty(t, st.Written())
	assert.Size(t, 3, st.Buffer)
}
