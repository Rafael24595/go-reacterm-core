package marker

// HistoryMeta defines the labels and separator glyphs used to format navigation history bars or breadcrumbs.
type HistoryMeta struct {
	// BackTag is the prefix label or symbol representing previous navigation history (e.g. "Back:", "←").
	BackTag string
	// NextTag is the prefix label or symbol representing forward navigation history (e.g. "Next:", "→").
	NextTag string
	// Separator is the inline delimiter string rendered between history entries (e.g. " | ", " / ").
	Separator string
}

// DefaultHistory defines a standard ASCII navigation history framing style.
var DefaultHistory = HistoryMeta{
	BackTag:   "Back:",
	NextTag:   "Next:",
	Separator: " | ",
}
