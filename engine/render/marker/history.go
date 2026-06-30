package marker

var DefaultHistory = HistoryMeta{
	BackTag:   "Back:",
	NextTag:   "Next:",
	Separator: " | ",
}

type HistoryMeta struct {
	BackTag   string
	NextTag   string
	Separator string
}
