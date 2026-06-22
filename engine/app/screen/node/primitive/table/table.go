package table

import (
	assert "github.com/Rafael24595/go-assert/assert/runtime"
	
	"github.com/Rafael24595/go-reacterm-core/engine/app/pager/predicate"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen"
	"github.com/Rafael24595/go-reacterm-core/engine/app/screen/keymap"
	"github.com/Rafael24595/go-reacterm-core/engine/app/state"
	"github.com/Rafael24595/go-reacterm-core/engine/app/viewmodel"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/cols"
	"github.com/Rafael24595/go-reacterm-core/engine/config/padding/rows"
	"github.com/Rafael24595/go-reacterm-core/engine/helper/math"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/decorator/inputline"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/drain"
	"github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/stream/pipeline/padding"
	"github.com/Rafael24595/go-reacterm-core/engine/model/hint"
	"github.com/Rafael24595/go-reacterm-core/engine/model/input"
	"github.com/Rafael24595/go-reacterm-core/engine/model/table"
	"github.com/Rafael24595/go-reacterm-core/engine/model/winsize"
	"github.com/Rafael24595/go-reacterm-core/engine/render/style"

	drawable_table "github.com/Rafael24595/go-reacterm-core/engine/layout/drawable/widget/table"
)

const Name = "table"

type MarshalFunc[T any] func(T) []table.Field

type Table[T any] struct {
	reference  string
	loaded     bool
	bindings   bindings
	definition definition
	action     *input.TableAction
	table      *table.Table
	cursor     *input.MatrixCursor
	positionY  style.VerticalPosition
	positionX  style.HorizontalPosition
}

func New[T any]() *Table[T] {
	return &Table[T]{
		reference:  Name,
		loaded:     false,
		bindings:   defaultBindings,
		definition: emptyDefinition(),
		action:     input.NewTableAction(),
		table:      table.NewTable(),
		cursor:     input.NewMatrixCursor(0, 0, false),
		positionY:  style.Middle,
		positionX:  style.Center,
	}
}

func (n *Table[T]) SetName(name string) *Table[T] {
	n.reference = name
	return n
}

func (n *Table[T]) WithWriteBindings(overrides *keymap.Bindings[CommandWrite]) *Table[T] {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.write = n.bindings.write.Overlay(overrides)
	return n
}

func (n *Table[T]) WithReadBindings(overrides *keymap.Bindings[CommandRead]) *Table[T] {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.bindings.read = n.bindings.read.Overlay(overrides)
	return n
}

func (n *Table[T]) EnableAction() *Table[T] {
	n.action.EnableMode = true
	return n
}

func (n *Table[T]) DisableAction() *Table[T] {
	n.action.EnableMode = false
	return n
}

func (n *Table[T]) SetActionHandler(handler input.TableActionHandler) *Table[T] {
	if n.loaded {
		assert.Unreachable(screen.MessageModified)
		return n
	}

	n.action.Handler = handler
	return n
}

func (n *Table[T]) SetPositionY(position style.VerticalPosition) *Table[T] {
	n.positionY = position
	return n
}

func (n *Table[T]) SetPositionX(position style.HorizontalPosition) *Table[T] {
	n.positionX = position
	return n
}

func (n *Table[T]) SetHeaders(headers ...string) *Table[T] {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	n.table = table.NewTable()
	n.table.SetHeaders(headers...)
	return n
}

func (n *Table[T]) AddItems(marshal MarshalFunc[T], items ...T) *Table[T] {
	if n.loaded {
		assert.Unreachable(screen.MessageNewElement)
		return n
	}

	rows := n.table.RowCount()

	for i, item := range items {
		index := rows + uint16(i)
		for _, field := range marshal(item) {
			n.table.SetCell(field.Header, index, field.Value)
		}
	}

	return n
}

func (n *Table[T]) ToNode() screen.Node {
	return screen.NewBuilder().
		Name(n.reference).
		NameToStack().
		Boot(n.boot).
		Keys(n.keys).
		Tick(n.tick).
		View(n.view).
		ToNode()
}

func (n *Table[T]) boot(uiState state.UIState) {
	if n.loaded {
		return
	}

	n.loaded = true

	n.loadFromStore(uiState)
	n.definition = definitionFromBindings(n.bindings)
}

func (n *Table[T]) loadFromStore(uiState state.UIState) {
	state, ok := KeyState.Get(
		uiState.Store,
		n.reference,
	)

	if !ok {
		return
	}

	n.cursor.Row = min(n.table.RowCount(), state.Row)
	n.cursor.Col = min(n.table.ColCount(), state.Col)
}

func (n *Table[T]) keys() screen.Definition {
	if !n.action.EnableMode {
		return screen.EmptyDefinition()
	}
	return n.definition.get(n.action.WriteMode)
}

func (n *Table[T]) tick(uiState *state.UIState, event screen.Event) screen.Result {
	uiState.Pager.ForceShow = true

	if !n.action.EnableMode {
		return screen.ResultFromUIState(uiState)
	}

	if !n.action.WriteMode {
		return n.tickRead(uiState, event)
	}
	return n.tickWrite(uiState, event)
}

func (n *Table[T]) tickWrite(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.write.Command(event.Key.Code) {
	case CmdWriteReadMode:
		n.action.WriteMode = false
		n.cursor.Show = n.action.WriteMode
	case CmdWriteMoveLeft:
		n.cursor.DecCol()
		n.tickToStore(uiState)
	case CmdWriteMoveRight:
		n.cursor.IncCol(
			math.SubClampZero(n.table.ColCount(), 1),
		)
		n.tickToStore(uiState)
	case CmdWriteMoveUp:
		n.cursor.DecRow()
		n.tickToStore(uiState)
	case CmdWriteMoveDown:
		n.cursor.IncRow(
			math.SubClampZero(n.table.RowCount(), 1),
		)
		n.tickToStore(uiState)
	}

	return screen.ResultFromUIState(uiState)
}

func (n *Table[T]) tickRead(uiState *state.UIState, event screen.Event) screen.Result {
	switch n.bindings.read.Command(event.Key.Code) {
	case CmdReadWriteMode:
		n.action.WriteMode = true
		n.cursor.Show = n.action.WriteMode
	}

	return screen.ResultFromUIState(uiState)
}

func (n *Table[T]) tickToStore(uiState *state.UIState) {
	tableState := State{
		Row: n.cursor.Row,
		Col: n.cursor.Col,
	}

	KeyState.Set(
		uiState.Store,
		n.reference,
		tableState,
	)
}

func (n *Table[T]) view(uiState state.UIState) viewmodel.ViewModel {
	vm := viewmodel.New()

	n.loadFromStore(uiState)

	table := drawable_table.UnitFromTable(*n.table, *n.cursor)

	position := padding.NewBuilder().
		Rows(
			hint.Maximize[winsize.Rows](),
			rows.WithPosition(n.positionY),
		).
		Cols(
			hint.Maximize[winsize.Cols](),
			cols.WithPosition(n.positionX),
		).
		ToUnit(table)

	vm.Kernel.Push(position)

	preficate := predicate.Page()
	if n.action.EnableMode && n.action.WriteMode {
		preficate = predicate.Focus()

		cell, _ := n.table.FindCellByCoords(n.cursor.Row, n.cursor.Col)

		vm.Footer.Push(
			inputline.Wrap(
				drain.UnitFromString(cell),
			),
		)
	}

	vm.Pager.SetPredicate(preficate)

	return *vm
}
