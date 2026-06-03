package table

import (
	"fmt"
	"testing"

	assert "github.com/Rafael24595/go-assert/assert/test"

	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/store"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"

	screen_test "github.com/Rafael24595/go-reacterm-core/test/engine/app/screen"
)

func TestTable_ToNode(t *testing.T) {
	node := New[int]().SetName("base").ToNode()

	screen_test.Helper_ToNode(t, node)
	assert.Equal(t, node.Name, "base")
}

func TestIndexMenu_Init(t *testing.T) {
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

	store.Push(
		uiState.Store,
		menu.reference,
		ArgTableState,
		State{
			Row: 1,
			Col: 3,
		},
	)

	node.Screen.Init(*uiState)

	assert.Equal(t, 1, menu.cursor.Row)
	assert.Equal(t, 3, menu.cursor.Col)
}

func TestTable_Stack(t *testing.T) {
	stack := New[int]().ToNode().Stack

	assert.True(t, stack.Has(Name))
}
