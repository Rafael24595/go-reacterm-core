package input

type CheckActionHandler = func()

func voidCheckHandler() {}

type CheckAction struct {
	WriteMode bool
	Handler   CheckActionHandler
}

func NewCheckAction(handler CheckActionHandler) *CheckAction {
	return &CheckAction{
		WriteMode: false,
		Handler:   handler,
	}
}

func EmptyCheckAction() *CheckAction {
	return NewCheckAction(voidCheckHandler)
}
