package table

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTable_ToNode(t *testing.T) {
	node := New[int]().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)
	assert.Equal(t, node.Name, "base")
}

func TestIndexMenu_Boot(t *testing.T) {
	menu := New[int]().
		SetHeaders("col_1", "col_2", "col_3").
		AddItems(
			func(i int) []table.Field {
				return []table.Field{
					{
						Header: fmt.Sprintf("col_%d", i),
						Value:  i,
					},
				}
			},
			1, 2, 3,
		)
	node := menu.ToNode()

	uiState := state.NewUIState()

	row := uint16(1)
	col := uint16(3)

	KeySync.Set(
		uiState.Store,
		menu.reference,
		Sync{
			Row: &row,
			Col: &col,
		},
	)

	node.Screen.Boot(*uiState)

	_, ok := KeySync.Get(uiState.Store, menu.reference)

	assert.False(t, ok)

	assert.Equal(t, row, menu.cursor.Row)
	assert.Equal(t, col, menu.cursor.Col)
}

func TestTable_Stack(t *testing.T) {
	stack := New[int]().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
