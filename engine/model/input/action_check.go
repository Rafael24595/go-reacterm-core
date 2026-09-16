package input

type CheckActionHandler = func()

func voidCheckHandler() {}

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

func DefaultCheckAction() *CheckAction {
	return newCheckAction(false, voidCheckHandler)
}

func (a *CheckAction) InEditMode() bool {
	return a.editMode
}

func (a *CheckAction) AsView() *CheckAction {
	a.editMode = false
	return a
}

func (a *CheckAction) AsEdit() *CheckAction {
	a.editMode = true
	return a
}

func (a *CheckAction) WithHandler(handler CheckActionHandler) *CheckAction {
	if handler == nil {
		return a
	}

	a.handler = handler
	return a
}

func (a CheckAction) Exec() {
	if a.handler != nil {
		a.handler()
	}
}
