package input

// TableActionHandler represents a callback signature executed with a matrix cursor context.
type TableActionHandler = func(MatrixCursor)

func voidTableHandler(_ MatrixCursor) {}

// TableAction encapsulates a matrix/table interaction callback alongside mode flags.
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

// DefaultTableAction returns a disabled no-op TableAction.
func DefaultTableAction() *TableAction {
	return newTableAction(false, false, voidTableHandler)
}

// IsNavigable reports whether navigation/interaction is enabled for this table action.
func (a *TableAction) IsNavigable() bool {
	return a.navigable
}

// EnableNavigation enables table navigation for this action.
func (a *TableAction) EnableNavigation() *TableAction {
	a.navigable = true
	return a
}

// DisableNavigation disables table navigation while preserving other settings.
func (a *TableAction) DisableNavigation() *TableAction {
	a.navigable = false
	return a
}

// InEditMode checks whether the action requires write permissions.
func (a *TableAction) InEditMode() bool {
	return a.editMode
}

// AsView sets the action mode to read-only.
func (a *TableAction) AsView() *TableAction {
	a.editMode = false
	return a
}

// AsEdit sets the action mode to write.
func (a *TableAction) AsEdit() *TableAction {
	a.editMode = true
	return a
}

// WithHandler updates the underlying callback handler and enables the action.
func (a *TableAction) WithHandler(handler TableActionHandler) *TableAction {
	if handler == nil {
		return a
	}

	a.navigable = true
	a.handler = handler

	return a
}

// Exec safely invokes the underlying handler if the action is enabled.
func (a *TableAction) Exec(cursor MatrixCursor) {
	if a.navigable && a.handler != nil {
		a.handler(cursor)
	}
}
