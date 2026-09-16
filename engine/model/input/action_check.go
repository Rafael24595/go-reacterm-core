package input

type CheckActionHandler = func()

func voidCheckHandler() {}

// CheckAction encapsulates an input callback alongside execution context flags.
type CheckAction struct {
	editMode bool
	handler  CheckActionHandler
}

func newCheckAction(
	editMode bool,
	handler CheckActionHandler,
) *CheckAction {
	return &CheckAction{
		editMode: editMode,
		handler:  handler,
	}
}

// DefaultCheckAction returns a no-op CheckAction.
func DefaultCheckAction() *CheckAction {
	return newCheckAction(false, voidCheckHandler)
}

// InEditMode checks whether the action requires write permissions/mode.
func (a *CheckAction) InEditMode() bool {
	return a.editMode
}

// AsView mutates and sets the action mode to read.
func (a *CheckAction) AsView() *CheckAction {
	a.editMode = false
	return a
}

// AsEdit mutates and sets the action mode to write.
func (a *CheckAction) AsEdit() *CheckAction {
	a.editMode = true
	return a
}

// WithHandler updates the underlying callback handler.
func (a *CheckAction) WithHandler(handler CheckActionHandler) *CheckAction {
	if handler == nil {
		return a
	}

	a.handler = handler
	return a
}

// Exec safely invokes the underlying handler.
func (a CheckAction) Exec() {
	if a.handler != nil {
		a.handler()
	}
}
