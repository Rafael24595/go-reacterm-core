package input

type TableActionHandler = func(MatrixCursor)

func voidTableHandler(_ MatrixCursor) {}

type TableAction struct {
	EnableMode bool
	WriteMode  bool
	Handler    TableActionHandler
}

func NewTableAction(handler ...TableActionHandler) *TableAction {
	enable := false
	handle := voidTableHandler
	if len(handler) > 0 {
		enable = true
		handle = handler[0]
	}

	return &TableAction{
		EnableMode: enable,
		WriteMode:  false,
		Handler:    handle,
	}
}

func (a *TableAction) SetHandler(handler TableActionHandler) *TableAction {
	a.EnableMode = true
	a.Handler = handler
	return a
}
