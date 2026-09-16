package input

type TableActionHandler = func(MatrixCursor)

func voidTableHandler(_ MatrixCursor) {}

type TableAction struct {
	navigable bool
	editMode  bool
	handler   TableActionHandler
}

func newTableAction(
	navigable bool,
	editMode bool,
	handler TableActionHandler,
) *TableAction {
	return &TableAction{
		navigable: navigable,
		editMode:  editMode,
		handler:   handler,
	}
}

func DefaultTableAction() *TableAction {
	return newTableAction(false, false, voidTableHandler)
}

func (a *TableAction) IsNavigable() bool {
	return a.navigable
}

func (a *TableAction) EnableNavigation() *TableAction {
	a.navigable = true
	return a
}

func (a *TableAction) DisableNavigation() *TableAction {
	a.navigable = false
	return a
}

func (a *TableAction) InEditMode() bool {
	return a.editMode
}

func (a *TableAction) AsView() *TableAction {
	a.editMode = false
	return a
}

func (a *TableAction) AsEdit() *TableAction {
	a.editMode = true
	return a
}

func (a *TableAction) WithHandler(handler TableActionHandler) *TableAction {
	if handler == nil {
		return a
	}

	a.navigable = true
	a.handler = handler

	return a
}

func (a *TableAction) Exec(cursor MatrixCursor) {
	if a.navigable && a.handler != nil {
		a.handler(cursor)
	}
}
